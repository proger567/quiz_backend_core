package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"quiz_backend_core/internal/config"
	"quiz_backend_core/internal/notifier"
	"quiz_backend_core/internal/service"
	"quiz_backend_core/internal/storage"
	transport "quiz_backend_core/internal/transport/http"
	"strings"
	"syscall"
	"time"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func LogLevel(lvl string) logrus.Level {
	switch strings.ToUpper(lvl) {
	case "INFO":
		return logrus.InfoLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "TRACE":
		return logrus.TraceLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "FATAL":
		return logrus.FatalLevel
	case "WARN":
		return logrus.WarnLevel
	default:
		panic("Not supported log level: " + lvl)
	}
}

func main() {
	// context
	mainCtx := context.Background()

	//config
	var cfg config.Config
	err := cfg.Parse()
	if err != nil {
		log.Fatal(err)
	}

	// DB
	databaseDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName)

	storages, err := storage.NewStorages(mainCtx, databaseDSN)
	if err != nil {
		log.Fatal(fmt.Errorf("Unable to connect to database: %v\n", err))
	}
	defer storages.Close()

	// logger
	logger := logrus.Logger{
		Out:   os.Stdout,
		Level: LogLevel(cfg.LogLevel),
		//ReportCaller: true,
		Formatter: &logrus.JSONFormatter{},
	}

	//notifier
	notifier, err := notifier.NewNotifier(cfg.NotifierHost+":"+cfg.NotifierPort, &logger)
	if err != nil {
		log.Fatal(fmt.Errorf("Unable to connect to notifier: %v\n", err))
	}
	defer notifier.Close()

	// metrics
	fieldKeys := []string{"method", "error"}
	requestCounter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api_test_generate",
		Subsystem: "login",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatencyMeter := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "api_test_generate",
		Subsystem: "login",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds",
	}, fieldKeys)

	// service
	deps := service.Deps{
		Storages:            storages,
		RequestCounter:      requestCounter,
		RequestLatencyMeter: requestLatencyMeter,
		Logger:              &logger,
		Notifier:            notifier,
	}

	s := service.NewServices(deps)

	// handler
	h := transport.MakeHTTPHandler(s, &logger, cfg.AppEnv == "development")

	addr := cfg.ListenAddr
	srv := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Error during ListenAndServ(). ", err)
			//os.Exit(1)
		}
	}()

	time.Sleep(1 * time.Second)
	logger.Info("Start service. Listen HTTP addr=", addr)

	<-done
	logger.Info("Server stopped. Signal ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server Exited error = ", err.Error())
	}
	logger.Info("Server Exited Properly")
}
