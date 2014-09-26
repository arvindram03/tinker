package main

import (
	"fmt"
	"code.google.com/p/go.net/websocket"
	"net/http"
	"tinkermaster/models"
	"log"
)

var events []models.Event

var regChannel chan *websocket.Conn
var unregChannel chan *websocket.Conn
var eventChannel chan models.Event


func run() {
	conns := make(map[*websocket.Conn]int)
	for {
		select {
		case conn := <-regChannel:
			log.Println("Registering ", conn)
			conns[conn] = 1
		case conn := <-unregChannel:
			log.Println("UnRegistering ", conn)
			conn.Close()
			delete(conns, conn)
		case event := <-eventChannel:
			events = append(events, event)
			for conn := range conns {
				log.Println("Sending ", conn)
				err := websocket.JSON.Send(conn, event)
				if err != nil {
					log.Println("Error sending event", err)
				}
			}
		}
	}
}


func Receiver(client *websocket.Conn) {
	regChannel <- client
	for {
		var rawEvent models.RawEvent
		err := websocket.JSON.Receive(client, &rawEvent)
		if err != nil {
			log.Println("err")
			unregChannel <- client
		}
		event := rawEvent.Parse()

		eventChannel <- event
		fmt.Println(events)
	}

}



func main() {
	regChannel = make(chan *websocket.Conn)
	unregChannel = make(chan *websocket.Conn)
	eventChannel = make(chan models.Event)
	http.Handle("/sync", websocket.Handler(Receiver))

	go run()
	err := http.ListenAndServe(":6055", nil)
	if err != nil {
		panic("Unable to listen and serve" + err.Error())
	}
}


