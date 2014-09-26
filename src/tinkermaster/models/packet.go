package models

import (
	"time"
	"log"
)
const (
	CREATE = "CREATE"
	DELETE = "DELETE"
	UPDATE = "UPDATE"
	ID = "_id"
)
type RawEvent struct {
	EventType		string	`json:"eventType"`
	Timestamp		string	`json:"timestamp"`
	Payload			interface {} `json:"payload"`
}


type Event struct {
	EventType		string	`json:"eventType"`
	Timestamp     	time.Time	`json:"timestamp"`
	Payload			interface {} `json:"payload"`
}

func (this RawEvent) Parse() Event {
	var event Event

	event.Payload = this.Payload
	event.EventType = this.EventType
	timestamp, err := time.Parse(time.RFC1123, this.Timestamp)
	if err != nil {
		log.Println("Unable to parse timestamp ", timestamp)
	}

	event.Timestamp = timestamp

	return event
}

