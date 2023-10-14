package controller

import (
	"log"
	"realtime-chat/translator"
	"sync"
	"time"

	"realtime-chat/model"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type EventHandler func(event model.Event, c *socketClient) error

type socketController struct {
	pongWait      time.Duration
	pingInterval  time.Duration
	activeClients ActiveClient
	translator    *translator.UTtrans
	sync.RWMutex
	handlers map[string]EventHandler
}

func (s *socketController) Serve(c *websocket.Conn) {
	log.Println("[Debug]: New Connection")
	upgraded := c.Locals("ws")

	if !upgraded.(bool) {
		log.Println(c.WriteJSON(fiber.Map{
			"error": s.translator.TranslateMessage(c.Locals("locale").(string), "ws", nil, nil),
		}))
		return
	}

	done := make(chan bool)
	client := newSocketClient(c, s, done)

	s.addClient(client)

	// Start client processes
	go client.readMessages()
	go client.writeMessages()
	<-done

}

// Add the client
func (s *socketController) addClient(client *socketClient) {
	s.Lock()
	defer s.Unlock()
	s.activeClients[client] = true
}

// Remove the client
func (s *socketController) removeClient(client *socketClient) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.activeClients[client]; ok {
		client.connection.Close()
		delete(s.activeClients, client)
	}

}

// set up all the available event handlers
func (s *socketController) setupEventHandlers() {
	s.handlers[model.SendMessageEvent] = func(e model.Event, c *socketClient) error {
		log.Println(e)
		return nil
	}

}

// look up the available event handler , and based on the messsage type , route to the approriate event handler
func (s *socketController) routeEvent(event model.Event, c *socketClient) error {
	if handler, ok := s.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return EvenNotSupportedError
	}

}

func NewSocketController(t *translator.UTtrans, pongWait time.Duration, pingInterval time.Duration) *socketController {
	c := &socketController{
		pongWait:      pongWait,
		pingInterval:  pingInterval,
		activeClients: make(ActiveClient),
		translator:    t,
		handlers:      make(map[string]EventHandler),
	}
	c.setupEventHandlers()
	return c
}
