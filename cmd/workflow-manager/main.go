package main

import (
	"context"
	"github.com/fidesy-pay/workflow-manager/internal/app"
	"github.com/fidesy-pay/workflow-manager/internal/config"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/completer"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/runners/registration"
	restorepassword "github.com/fidesy-pay/workflow-manager/internal/pkg/runners/restore_password"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/storage"
	auth_service "github.com/fidesy-pay/workflow-manager/pkg/auth-service"
	clients_service "github.com/fidesy-pay/workflow-manager/pkg/clients-service"
	crypto_service "github.com/fidesy-pay/workflow-manager/pkg/crypto-service"
	email_service "github.com/fidesy-pay/workflow-manager/pkg/email-service"
	"github.com/fidesy/sdk/common/grpc"
	"github.com/fidesy/sdk/common/logger"
	"github.com/fidesy/sdk/common/postgres"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	defer cancel()

	err := config.Init()
	if err != nil {
		log.Fatalf("config.Init: %v", err)
	}

	server, err := grpc.NewServer(
		grpc.WithPort(os.Getenv("GRPC_PORT")),
		grpc.WithMetricsPort(os.Getenv("METRICS_PORT")),
		grpc.WithDomainNameService(ctx, "domain-name-service:10000"),
		grpc.WithGraylog("graylog:5555"),
		grpc.WithTracer("http://jaeger:14268/api/traces"),
	)
	if err != nil {
		log.Fatalf("grpc.NewServer: %v", err)
	}

	pool, err := postgres.Connect(ctx, os.Getenv("PG_DSN"))
	if err != nil {
		logger.Fatalf("postgres.Connect: %v", err)
	}

	storage := storage.New(pool)

	emailClient, err := grpc.NewClient[email_service.EmailServiceClient](
		ctx,
		email_service.NewEmailServiceClient,
		"rpc:///email-service",
	)
	if err != nil {
		log.Fatalf("NewEmailClient: %v", err)
	}

	authClient, err := grpc.NewClient[auth_service.AuthServiceClient](
		ctx,
		auth_service.NewAuthServiceClient,
		"rpc:///auth-service",
	)
	if err != nil {
		log.Fatalf("NewAuthClient: %v", err)
	}

	clientsClient, err := grpc.NewClient[clients_service.ClientsServiceClient](
		ctx,
		clients_service.NewClientsServiceClient,
		"rpc:///clients-service",
	)
	if err != nil {
		log.Fatalf("NewClientsClient: %v", err)
	}

	cryptoClient, err := grpc.NewClient[crypto_service.CryptoServiceClient](
		ctx,
		crypto_service.NewCryptoServiceClient,
		"rpc:///crypto-service",
	)
	if err != nil {
		log.Fatalf("NewCryptoClient: %v", err)
	}

	registrationRunner := registration.NewRunner(
		storage,
		emailClient,
		authClient,
		clientsClient,
		cryptoClient,
	)

	restorePasswordRunner := restorepassword.NewRunner(
		storage,
		emailClient,
		authClient,
		clientsClient,
	)

	errGroup := errgroup.Group{}

	runner := completer.New(storage, registrationRunner, restorePasswordRunner)
	//errGroup.Go(func() error {
	//	if err = compl.Complete(ctx); err != nil {
	//		return fmt.Errorf("completer.Complete: %w", err)
	//	}
	//
	//	return nil
	//})

	impl := app.New(storage, runner)

	errGroup.Go(func() error {
		if err = server.Run(ctx, impl); err != nil {
			log.Fatalf("app.Run: %v", err)
		}

		return nil
	})

	if err = errGroup.Wait(); err != nil {
		logger.Fatalf("%w", err)
	}
}
