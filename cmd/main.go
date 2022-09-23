package main

import (
	"context"
	"fmt"

	"project/internal/adapter/exchange"
	"project/internal/adapter/repository"
	"project/internal/api"
	"project/internal/config"
	"project/internal/controller"
	"project/internal/scheduler"

	"github.com/gin-gonic/gin"
	pgx "github.com/jackc/pgx/v4"
	"github.com/robfig/cron/v3"

	_ "project/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/sirupsen/logrus"
)

// @title Exchanger
// @version 1.0
// @description Swagger API for Golang Project Exchanger

// @contact.name   Ilya
// @contact.email    biv_1998@mail.ru

// @host      localhost:8080
// @BasePath /
func main() {
	// Инициализация конфигурации
	config, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalln("init config: ", err)
	}

	// Подключение к базе данных
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		logrus.Fatalln("connect to DB", err)
	}
	defer db.Close(context.Background())

	repository, err := repository.NewRepository(db)
	if err != nil {
		logrus.Fatalln("init repository: ", err)
	}
	exchange, err := exchange.NewExchangeClient(config.Apikey)
	if err != nil {
		logrus.Fatalln("init exchange: ", err)
	}
	controller, err := controller.NewController(repository, exchange)
	if err != nil {
		logrus.Fatalln("init controller: ", err)
	}
	api, err := api.NewAPI(controller)
	if err != nil {
		logrus.Fatalln("init api: ", err)
	}

	// Настройка и запуск планировщика
	c := cron.New()
	s, err := scheduler.NewScheduler(controller, c, config.Schedule)
	if err != nil {
		logrus.Fatalln("init scheduler: ", err)
	} else {
		s.Start()
	}

	// Настройка роутов и запуск сервера
	r := gin.Default()
	r.GET("/currencies", api.CurrenciesHandler)
	r.GET("/rates", api.RateHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err = r.Run(":8080")
	if err != nil {
		logrus.Fatalln("router run: ", err)
	}
}
