package runner

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
	"user-subscriptions/internal/app"
	"user-subscriptions/internal/app/command"
	"user-subscriptions/internal/app/query"
	"user-subscriptions/internal/config"
	"user-subscriptions/internal/port/http"
	"user-subscriptions/internal/repository"
	"user-subscriptions/internal/server"
)

func Start(configDir, migrationDir string) {
	cfg := newConfig(configDir)

	db := initDB(cfg, migrationDir)
	defer db.Close()
	application := newApplication(db)

	serv := startServer(cfg, application)
	shutdown(serv)
}

func newConfig(configDir string) *config.Config {
	cfg, err := config.New(configDir)
	if err != nil {
		log.Panicln(err)
	}

	return cfg
}

func initDB(cfg *config.Config, migrationDir string) *sql.DB {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.Name)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Panicln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Panicln(err)
	}

	err = upMigration(db, migrationDir)
	if err != nil {
		logrus.Panicln(err)
	}

	return db
}

func upMigration(db *sql.DB, migrationDir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}
	err = goose.Status(db, migrationDir)
	if err != nil {
		return err
	}
	return goose.Up(db, migrationDir)
}

func newApplication(db *sql.DB) app.Application {
	subsRepository := repository.NewSubsRepository(db)
	return app.Application{
		Queries: app.Queries{
			GetSub:          query.NewGetSubHandler(subsRepository),
			GetSubsByUserId: query.NewGetSubUsrIdHandler(subsRepository),
			GetSubsPrice:    query.NewGetSubsUserServiceDate(subsRepository),
		},
		Commands: app.Commands{
			AddSubUserCommand: command.NewAddSubHandler(subsRepository),
			RemoveSubCommand:  command.NewRemoveSubHandler(subsRepository),
			UpdateSubCommand:  command.NewUpdateSubHandler(subsRepository),
		},
	}
}

func startServer(cfg *config.Config, application app.Application) *server.Server {
	logrus.Info(fmt.Sprintf("Starting HTTP server on address: %s", cfg.HTTP.Port))
	httpServer := server.New(cfg, http.NewHandler(application))

	go func() {
		err := httpServer.Run()
		if err != nil {
			log.Fatalf("HTTP server stopped, %s\n", err)
		}
	}()

	return httpServer
}

func shutdown(server *server.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	logrus.Println("Shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Fatal(err)
	}
}
