package controller

import (
	"encoding/json"
	"log"
	"realtime-chat/model"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type ActiveClient map[*socketClient]bool

type socketClient struct {
	connection *websocket.Conn
	controller *socketController
	// Avoid concurrent write on the websocket
	egress chan []byte
	done   chan<- bool
}

// Read the messages send by client
func (s *socketClient) readMessages() {

	defer func() {
		s.done <- true
		s.controller.removeClient(s)
	}()

	s.connection.SetReadLimit(512)

	// the initail deadline timer
	if err := s.connection.SetReadDeadline(time.Now().Add(s.controller.env.PongWait)); err != nil {
		log.Println("Initial deadline breach", err)
		return
	}

	s.connection.SetPongHandler(s.pongHandler)

	var (
		msg []byte
		err error
	)

	for {

		if _, msg, err = s.connection.ReadMessage(); err != nil {

			// If not coming from client or server
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			log.Printf("Read error : %v", err)

			break
		}

		var request model.Event
		if err := json.Unmarshal(msg, &request); err != nil {
			log.Printf("error marshalling message: %v", err)
			err = s.writeError(err)
			if err != nil {
				break
			}
		}

		if err := s.controller.routeEvent(request, s); err != nil {
			log.Println("Error handeling Message:", err)
			break

		}

	}
}

// Write messages
func (s *socketClient) writeMessages() {

	ticker := time.NewTicker(s.controller.env.PingInterval)

	defer func() {
		ticker.Stop()
		s.controller.removeClient(s)
	}()

	for {
		select {
		case message, ok := <-s.egress:
			if !ok {
				if err := s.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection Closed:", err)

				}
				return

			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			if err := s.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}
			log.Println("sent message")
		case <-ticker.C:
			if err := s.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("pingmsg", err)
				return

			}
		}

	}

}

func (s *socketClient) writeError(err error) error {

	payload := model.Event{Type: model.ErrorMessageEvent, Payload: json.RawMessage(err.Error())}
	b, err := json.Marshal(&payload)
	if err != nil {

		log.Println("Marshal error", err)
		return err
	}

	if err := s.connection.WriteMessage(websocket.TextMessage, b); err != nil {
		log.Println("Connection Closed:", err)
		return err
	}
	return nil

}

// Pong message handler , by continute set the deadline of message
func (s *socketClient) pongHandler(msg string) error {
	log.Println("pong")
	return s.connection.SetReadDeadline(time.Now().Add(s.controller.env.PongWait))
}

// Return the new socket client of the given manager and underlying connection
func newSocketClient(conn *websocket.Conn, controller *socketController, done chan<- bool) *socketClient {
	return &socketClient{
		connection: conn,
		controller: controller,
		egress:     make(chan []byte),
		done:       done,
	}

}
