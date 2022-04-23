package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {

	app := fiber.New()
	go websocketHub()
	app.Use("/ws", WebsocketUpgrade)
	app.Get("/ws/chat", websocket.New(WebSocketRun))

	app.Static("/", "./web")

	// ws://localhost:3000/ws
	log.Fatal(app.Listen(":3000"))
}
