package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/tanishk-deore/url-shortner/handlers"
	"github.com/tanishk-deore/url-shortner/logger"
	"github.com/tanishk-deore/url-shortner/service"
	"github.com/tanishk-deore/url-shortner/store"
)

func main() {

	//this is default code

	// memoryStore:= store.NewMemoryStore()

	// shortnerService := service.NewShortnerService(
	// 	memoryStore,
	// )

	// urlHandler := handlers.NewUrlHandler(
	// 	shortnerService,
	// )

	// http.HandleFunc(
	// "/shorten",
	// urlHandler.ShortenUrl,
	// )

	// http.HandleFunc(
	// 	"/",
	// 	urlHandler.RedirectURL,
	// )

	// fmt.Println("Server Started on:8080")

	// err:= http.ListenAndServe(
	// 	":8080",
	// 	nil,
	// )

	// if err!= nil {
	// 	panic(err)
	// }

	//addding the postgres db
	//will check if url is present in the env
	//if it is it will use it otherwise by default it will be the same as above
	//storing the database in runtime in memory

	//update
	//changed it to default postgres here removed the local memory store
	//for multiple

	var shortnerService *service.ShortnerService

	dsn := os.Getenv("DATABASE_URL")
	logger.Info("Loading db url from env...")
	if dsn == "" {
		panic("Database Url must be in set")
	}

	logger.Info("Db url found and initializing postgrestore")

	pgStore, err := store.NewPostgresStore(dsn)

	if err != nil {
		panic(err)
	}

	logger.Info("Db pool Successfully created")
	//now we need to setup both redis and postgres

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisCache, err := store.NewRedisCache(redisAddr)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := redisCache.Close(); err != nil {
			logger.Error("Error closing Cache : %v", err)
		}
	}()

	shortnerService = service.NewShortnerService(pgStore, redisCache)

	// shortnerService = service.NewShortnerService(pgStore)

	urlHandler := handlers.NewUrlHandler(shortnerService)

	http.HandleFunc("/shorten", urlHandler.ShortenUrl)

	http.HandleFunc("/", urlHandler.RedirectURL)

	//updated the code here to check the por tnumber which will be passed as env
	//erailer it was by deafult at 8080
	//if no port number is provided then it will take it as 8080

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	logger.Info("Starting server one port: %s", port)

	err = http.ListenAndServe(":"+port, nil)

	if err != nil {
		logger.Error("Server Error: %v", err)
		panic(err)
	}

	fmt.Println("Server started on:" + port)

}
