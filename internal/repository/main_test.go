package repository

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/distuurbia/tradingService/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
	"github.com/sirupsen/logrus"
)

const (
	pgUsername = "postgresql"
	pgPassword = "postgresql"
	pgDB       = "positionsDB"
	pgPort     = "5432/tcp"
)

var (
	r            *TradingServiceRepository
	testPosition = model.Position{
		PositionID:  uuid.New(),
		ProfileID:   uuid.New(),
		MoneyAmount: 300,
		StopLoss:    10,
		TakeProfit:  200,
		ShareName:   "Apple",
		Vector:      "long",
	}
)

func SetupTestPostgres() (*pgxpool.Pool, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %w", err)
	}
	resource, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_USER=" + pgUsername,
		"POSTGRES_PASSWORD=" + pgPassword,
		"POSTGRES_DB=" + pgDB})
	if err != nil {
		logrus.Fatalf("can't start postgres container: %s", err)
	}
	cmd := exec.Command(
		"flyway",
		fmt.Sprintf("-user=%s", pgUsername),
		fmt.Sprintf("-password=%s", pgPassword),
		"-locations=filesystem:../../migrations",
		fmt.Sprintf("-url=jdbc:postgresql://%s:%s/%s", "localhost", resource.GetPort(pgPort), pgDB), "-connectRetries=10",
		"migrate",
	)
	err = cmd.Run()
	if err != nil {
		logrus.Fatalf("can't run migration: %s", err)
	}
	dbURL := fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s", pgUsername, pgPassword, resource.GetPort(pgPort), pgDB)
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logrus.Fatalf("can't parse config: %s", err)
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		logrus.Fatalf("can't connect to postgtres: %s", err)
	}
	cleanup := func() {
		dbpool.Close()
		pool.Purge(resource)
	}
	return dbpool, cleanup, nil
}

func TestMain(m *testing.M) {
	pool, cleanup, err := SetupTestPostgres()
	if err != nil {
		logrus.Fatalf("failed to construct the pool -> %v", err)
		cleanup()
		os.Exit(1)
	}
	r = NewTradingServiceRepository(pool)
	exitCode := m.Run()
	cleanup()
	os.Exit(exitCode)
}
