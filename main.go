package main

import (
	"boom/server"
	"log"
	"net/http"
)

func main() {

	m := server.NewServer()

	log.Println("🔥 Сервер жгёт!")
	log.Fatal(http.ListenAndServe(":8081", m))
}
