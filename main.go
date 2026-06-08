package main

import (
	"fmt"
	"net/http"
	"os"

	"url-shortner/handlers"
	"url-shortner/service"
	"url-shortner/store"
)

func main (){

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

	var shortnerService *service.ShortnerService

    if dsn := os.Getenv("DATABASE_URL"); dsn != "" {


        pgStore, err := store.NewPostgresStore(dsn)


        if err != nil {
            panic(err)
        }

        shortnerService = service.NewShortnerService(pgStore)
    } else {

        memoryStore := store.NewMemoryStore()

        shortnerService = service.NewShortnerService(memoryStore)
    }

    urlHandler := handlers.NewUrlHandler(shortnerService)


    http.HandleFunc("/shorten", urlHandler.ShortenUrl)



    http.HandleFunc("/", urlHandler.RedirectURL)

    fmt.Println("Server Started on:8080")


    err := http.ListenAndServe(":8080", nil)

    if err != nil {
        panic(err)
    }



}