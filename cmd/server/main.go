package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"github.com/nixpig/syringe.sh/internal/auth"
	"github.com/nixpig/syringe.sh/internal/database"
	"github.com/nixpig/syringe.sh/internal/middleware"
	"github.com/nixpig/syringe.sh/pkg/validation"
	"github.com/rs/zerolog"
)

func main() {
	// -- LOGGING
	log := zerolog.
		New(os.Stdout).
		Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02T15:04:05.999Z07:00",
		}).With().Timestamp().Logger()

	// -- ENV
	log.Info().Msg("loading environment")
	if err := godotenv.Load(".env"); err != nil {
		log.Error().Err(err).Msg("failed to load '.env' file")
		os.Exit(1)
	}

	// -- DATABASE
	log.Info().Msg("connecting to database")
	appDB, err := database.Connection(
		os.Getenv("DB_FILENAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		os.Exit(1)
	}

	defer appDB.Close()

	// -- RUN DB MIGRATION
	log.Info().Msg("running database migrations")

	migrator, err := database.NewMigration(appDB, "migrations/app")
	if err != nil {
		log.Error().Err(err).Msg("failed to create migration")
		os.Exit(1)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error().Err(err).Msg("failed to run migration")
		os.Exit(1)
	}

	// -- DEPENDENCY CONSTRUCTION
	log.Info().Msg("building app components")
	validate := validation.New()
	authStore := auth.NewSqliteAuthStore(appDB)
	authService := auth.NewAuthService(authStore, validate)

	// -- SERVER
	sshServer := newServer(
		&log,
		[]wish.Middleware{
			middleware.NewMiddlewareCommand(&log, appDB, validate),
			middleware.NewMiddlewareAuth(&log, authService),
			middleware.NewMiddlewareLogging(&log),
		},
		time.Duration(time.Second*30),
		os.Getenv("HOST_KEY_PATH"),
	)

	if err := sshServer.Start(
		os.Getenv("APP_HOST"),
		os.Getenv("APP_PORT"),
	); err != nil {
		log.Error().Err(err).Msg("failed to start ssh server")
		os.Exit(1)
	}
}

type Server struct {
	logger      *zerolog.Logger
	middleware  []wish.Middleware
	timeout     time.Duration
	hostKeyPath string
}

func newServer(
	logger *zerolog.Logger,
	middleware []wish.Middleware,
	timeout time.Duration,
	hostKeyPath string,
) Server {
	return Server{
		logger:      logger,
		middleware:  middleware,
		timeout:     timeout,
		hostKeyPath: hostKeyPath,
	}
}

func (s Server) Start(host, port string) error {
	server, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(s.hostKeyPath),
		wish.WithMaxTimeout(s.timeout),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			allowed := key.Type() == "ssh-rsa"
			return allowed
		}),
		wish.WithMiddleware(
			s.middleware...,
		),
	)
	if err != nil {
		return err
	}

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s.logger.Info().
		Str("host", host).
		Str("port", port).
		Msg("starting server")

	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			s.logger.Error().Err(err).Msg("failed to start server")
			done <- nil
		}
	}()

	s.logger.Info().Msg("server started")

	<-done

	s.logger.Info().Msg("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		s.logger.Error().Err(err).Msg("failed to stop server")
		return err
	}

	return nil
}
