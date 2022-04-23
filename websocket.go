package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type client struct {
	message string
}

// tüm clientleri tutan kanal
var clients = make(map[*websocket.Conn]client)

// tüm gelen mesajlarının broadcast edilmesi için
var brodcast = make(chan client)

// kayıtlı olan tüm clientleri tutan kanal
var register = make(chan *websocket.Conn)

//kayıtsız olan clientleri tutan kanal
var unregister = make(chan *websocket.Conn)

//clientler arası iletişim fonksiyonu
func websocketHub() {

	for { //sonsuz döngü
		select { // case select koşulu
		case connetion := <-register: // register durumunda gelen clientleri kanala ekliyoruz

			clients[connetion] = client{}

			log.Println("connected: ", len(clients))
		case wsup := <-brodcast: // broadcast durumunda gelen mesajları broadcast ediyoruz
			//tüm clientlere mesaj gönderiyoruz
			for client := range clients { // clients kanalindaki tüm clientleri döngüye sokuyoruz
				err := client.WriteMessage(websocket.TextMessage, []byte(wsup.message)) //cliente mesaj gönderiyoruz
				if err != nil {
					log.Println(err)

				}
			}
		case connetion := <-unregister: //unregister durumunda gelen clientleri kanalamızden çıkartıyoruz
			delete(clients, connetion)
			for client := range clients {
				client.WriteMessage(websocket.TextMessage, []byte("bir kullanıcı çıkış yaptı"))
			}
		}

	}

}

//Bu websocket upgrade fonksiyonu alt bir fonksiyon
func WebsocketUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()

	}
	return fiber.ErrUpgradeRequired
}

//Bu websocket gelen data fonksiyonu
func WebSocketRun(c *websocket.Conn) {

	//hoşgeldin mesajını cliente gönderiyoruz
	c.WriteMessage(websocket.TextMessage, []byte("Mr.Robot: Are you one or zero?"))
	register <- c //register listemize ekliyoruz
	for {         //sonsuz döngü bağlantı boyunca
		messageType, message, err := c.ReadMessage() //clientden gelen mesajı okuyoruz

		if err != nil { //hata alınırsa
			log.Println(err)
			unregister <- c //unregister listemizden çıkartıyoruz
			break
		}
		if messageType == websocket.TextMessage { //mesaj tipi text ise
			brodcast <- client{message: string(message)} //broadcast listemize ekliyoruz

		}
	}

}
