package main

import (
	"log"
	"os"

	todo "github.com/ch0c0-msk/example-todo-app"
	"github.com/ch0c0-msk/example-todo-app/pkg/handler"
	"github.com/ch0c0-msk/example-todo-app/pkg/repository"
	"github.com/ch0c0-msk/example-todo-app/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func initConfig() error {
	viper.SetConfigFile("configs/config.yml")
	return viper.ReadInConfig()
}

func init() {
	if err := initConfig(); err != nil {
		log.Fatalf("initizialtion error: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
}

func main() {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	handler := handler.NewHandler(service)

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
		log.Fatalf("running server error: %s", err.Error())
	}
}
