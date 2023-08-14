// Package main contains main goroutine and connections to other microservices
package main

import (
	"context"
	"fmt"

	"github.com/caarlos0/env"
	priceProtocol "github.com/distuurbia/PriceService/protocol/price"
	"github.com/distuurbia/tradingService/internal/config"
	"github.com/distuurbia/tradingService/internal/model"
	"github.com/distuurbia/tradingService/internal/repository"
	"github.com/distuurbia/tradingService/internal/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createPriceServicelientConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	priceServiceConn, err := createPriceServicelientConnection()
	if err != nil {
		logrus.Errorf("main -> : %v", err)
	}
	defer func() {
		errConnClose := priceServiceConn.Close()
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

	tradingRps := repository.NewTradingServiceRepository(pool)
	client := priceProtocol.NewPriceServiceServiceClient(priceServiceConn)
	priceRps := repository.NewPriceServiceRepository(client)
	s := service.NewTradingServiceService(priceRps, tradingRps)
	testProfileID := uuid.New()
	testPositionID := uuid.New()
	go s.SendSharesToProfiles(context.Background())
	s.AddInfoToManager(testProfileID, "Coca-Cola")
	if err != nil {
		logrus.Errorf("main -> : %v", err)
	}
	// go func() {
	// 	time.Sleep(2 *time.Second)
	// 	s.ClosePosition(context.Background(), testProfileID, testPositionID, "Coca-Cola")
	// }()
	amount, err := s.OpenPosition(context.Background(), &model.Position{
		ProfileID:   testProfileID,
		PositionID:  testPositionID,
		ShareName:   "Coca-Cola",
		MoneyAmount: 5000,
		TakeProfit:  125,
		StopLoss:    120,
		Vector:      "long",
	})
	if err != nil {
		logrus.Errorf("main -> : %v", err)
	}
	fmt.Println(amount)
}
