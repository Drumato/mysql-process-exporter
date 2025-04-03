package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"

	"database/sql"

	"github.com/Drumato/mysql-process-exporter/metrics"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultInitialiMySQLConnectRetryCount = 10
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	// Setup
	e := echo.New()
	reg := metrics.InitializeMetrics()
	e.GET("/healthz", func(c echo.Context) error {
		// Simulate a health check
		return c.JSON(http.StatusOK, "OK")
	})

	mysqlConfig := mysql.Config{
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		Net:    "tcp",
		Addr:   net.JoinHostPort(os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT")),
	}

	mysqlConfig.FormatDSN()
	var dbConn *sql.DB
	for i := 0; i < defaultInitialiMySQLConnectRetryCount; i++ {
		db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
		if err != nil {
			slog.ErrorContext(
				context.Background(),
				"failed to open mysql connection",
				slog.String("error", err.Error()),
				slog.Int("attempt", i+1),
				slog.Int("max_attempts", defaultInitialiMySQLConnectRetryCount), slog.Int("interval_seconds", 5),
			)
			time.Sleep(5 * time.Second)
			continue
		}
		if err := db.Ping(); err != nil {
			slog.ErrorContext(
				context.Background(),
				"failed to ping mysql",
				slog.String("error", err.Error()),
				slog.Int("attempt", i+1),
				slog.Int("max_attempts", defaultInitialiMySQLConnectRetryCount), slog.Int("interval_seconds", 5),
			)
			time.Sleep(5 * time.Second)
			continue
		}
		dbConn = db
		break
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			slog.ErrorContext(context.Background(), "failed to close mysql connection", slog.String("error", err.Error()))
		}
	}()

	e.Use(metrics.OndemandUpdateMetricsMiddleware(logger, mysqlConfig.Addr, dbConn))
	e.GET("/metrics", echo.WrapHandler(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	port, ok := os.LookupEnv("MYSQL_PROCESS_EXPORTER_PORT")
	if !ok {
		port = "8080"
	}

	defer stop()
	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "shutting down the server", slog.String("error", err.Error()))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
