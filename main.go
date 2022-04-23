package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/websocket/v2"
)

//go:embed web/*
var embedDirStatic embed.FS

func main() {

	app := fiber.New()
	go websocketHub()
	app.Use("/ws", WebsocketUpgrade)
	app.Get("/ws/chat", websocket.New(WebSocketRun))
	app.Use("/static", filesystem.New(filesystem.Config{
		Root:   http.FS(embedDirStatic),
		Browse: true,
	}))

	// ws://localhost:3000/ws
	log.Fatal(app.Listen(":3000"))
}
