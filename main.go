package main

import (
	"fmt"
	"net/http"

	"url-shortner/handlers"
	"url-shortner/service"
	"url-shortner/store"
)

func main (){

	memoryStore:= store.NewMemoryStore()

	shortnerService := service.NewShortnerService(
		memoryStore,
	)

	urlHandler := handlers.NewUrlHandler(
		shortnerService,
	)

	http.HandleFunc(
	"/shorten",
	urlHandler.ShortenUrl,
	)


	http.HandleFunc(
		"/",
		urlHandler.RedirectURL,
	)


	fmt.Println("Server Started on:8080")

	err:= http.ListenAndServe(
		":8080",
		nil,
	)

	if err!= nil {
		panic(err)
	}



}