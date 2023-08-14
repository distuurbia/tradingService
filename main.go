// Package main contains main goroutine and connections to other microservices
package main

import (
	"context"
	"fmt"
	"net"

	"github.com/caarlos0/env"
	priceProtocol "github.com/distuurbia/PriceService/protocol/price"
	balanceProtocol "github.com/distuurbia/balance/protocol/balance"
	"github.com/distuurbia/tradingService/internal/config"
	"github.com/distuurbia/tradingService/internal/handler"
	protocol "github.com/distuurbia/tradingService/protocol/trading"
	"github.com/go-playground/validator"

	"github.com/distuurbia/tradingService/internal/repository"
	"github.com/distuurbia/tradingService/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createPriceServiceClientConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func createBalanceClientConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func connectPostgres(cfg *config.Config) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(cfg.PostgresPath)
	if err != nil {
		return nil, fmt.Errorf("error in method pgxpool.ParseConfig: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		return nil, fmt.Errorf("error in method pgxpool.NewWithConfig: %v", err)
	}
	return pool, nil
}

func main() {
	priceServiceConn, err := createPriceServiceClientConnection()
	if err != nil {
		logrus.Errorf("main -> : %v", err)
	}
	defer func() {
		errConnClose := priceServiceConn.Close()
		if errConnClose != nil {
			logrus.Fatalf("main -> : %v", errConnClose)
		}
	}()

	balanceConn, err := createBalanceClientConnection()
	if err != nil {
		logrus.Errorf("main -> : %v", err)
	}
	defer func() {
		errConnClose := balanceConn.Close()
		if errConnClose != nil {
			logrus.Fatalf("main -> : %v", errConnClose)
		}
	}()

	var cfg config.Config
	if errCfg := env.Parse(&cfg); errCfg != nil {
		logrus.Fatalf("main -> %v", errCfg)
	}
	pool, errPgx := connectPostgres(&cfg)
	if errPgx != nil {
		logrus.Fatalf("main -> %v", errPgx)
	}

	priceServiceClient := priceProtocol.NewPriceServiceServiceClient(priceServiceConn)
	balanceClient := balanceProtocol.NewBalanceServiceClient(balanceConn)

	tradingRps := repository.NewTradingServiceRepository(pool)
	priceRps := repository.NewPriceServiceRepository(priceServiceClient, &cfg)
	balanceRps := repository.NewBalanceRepository(balanceClient)

	s := service.NewTradingServiceService(priceRps, balanceRps, tradingRps)

	validate := validator.New()
	h := handler.NewTradingServiceHandler(s, validate, &cfg)

	go s.SendSharesToProfiles(context.Background())

	lis, err := net.Listen("tcp", "localhost:8086")
	if err != nil {
		logrus.Fatalf("main -> %v", err)
	}

	serverRegistrar := grpc.NewServer()
	protocol.RegisterTradingServiceServiceServer(serverRegistrar, h)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		logrus.Fatalf("main -> %v", err)
	}
}
