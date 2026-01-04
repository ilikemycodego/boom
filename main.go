package main

import (
	"boom/db"
	"boom/server"
	"log"
	"net/http"
)

func main() {
	// 1. Открываем базу один раз при старте
	database := db.OpenDB("data/data.db")
	defer database.Close()

	// 2. Передаем её в сервер
	m := server.NewServer(database)

	log.Println("🔥 Сервер жгёт!")
	log.Fatal(http.ListenAndServe(":8081", m))
}
