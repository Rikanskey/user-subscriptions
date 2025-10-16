package runner

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
	"user-subscriptions/internal/app"
	"user-subscriptions/internal/app/command"
	"user-subscriptions/internal/app/query"
	"user-subscriptions/internal/config"
	"user-subscriptions/internal/port/http"
	"user-subscriptions/internal/repository"
	"user-subscriptions/internal/server"
)

func Start(configDir string) {
	cfg := newConfig(configDir)
	db := initDB(cfg)
	application := newApplication(db)
	startServer(cfg, application)
}

func newConfig(configDir string) *config.Config {
	cfg, err := config.New(configDir)
	if err != nil {
		log.Panicln(err)
	}

	return cfg
}

func initDB(cfg *config.Config) *sql.DB {
	dbInfo := fmt.Sprintf("postgresql://%s:%s@postgres/%s?sslmode=disable",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Name)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Panicln(err)
	}
	err = db.Ping()
	if err != nil {
		logrus.WithError(err).Println("Unable to connect to database")
	}
	return db
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

func startServer(cfg *config.Config, application app.Application) {
	logrus.Info(fmt.Sprintf("Starting HTTP server on address: %s", cfg.HTTP.Port))
	httpServer := server.New(cfg, http.NewHandler(application))

	err := httpServer.Run()

	log.Panicln("HTTP server stopped, ", err)
}
