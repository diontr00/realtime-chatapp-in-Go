package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"realtime-chat/config"
	"realtime-chat/internal/otp"
	"realtime-chat/translator"
	"strings"
	"sync"

	"realtime-chat/model"

	"github.com/go-playground/validator"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var authDB = map[string]string{"admin": "admin", "notadmin": "123"}

type EventHandler func(event model.Event, c *socketClient) error

type socketController struct {
	env           config.SocketEnv
	activeClients ActiveClient
	translator    *translator.UTtrans
	validator     *validator.Validate
	sync.RWMutex
	handlers map[string]EventHandler
	otps     otp.RetentionMap
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

	otp := c.Query("otp")
	if otp == "" || !s.otps.VerifyOTP(otp) {
		payload := model.Event{Type: model.ErrorMessageEvent, Payload: json.RawMessage("unauthorized")}

		b, err := json.Marshal(&payload)
		if err != nil {
			c.Close()
			return
		}
		if err := c.WriteMessage(websocket.TextMessage, b); err != nil {
			log.Println("Send error", err)

			c.Close()
			return
		}
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
	s.handlers[model.SendMessageEvent] = broadCast

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

// Authentication Handler
func (s *socketController) LoginHandler(c *fiber.Ctx) error {
	var req model.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Println(err)

		return c.Status(fiber.StatusBadRequest).JSON(model.UserLoginResponse{
			Status: "Bad Request",
		})

	}
	locale := c.Locals("locale").(string)

	if err := s.translator.ValidateRequest(locale, s.validator, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.UserLoginResponse{
			Status: "Bad Request",
			Error:  strings.Join(err, "\n"),
		})
	}

	if password, ok := authDB[req.Username]; !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(model.UserLoginResponse{
			Status: "Unauthorized",
		})
	} else if password != req.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(model.UserLoginResponse{
			Status: "Unauthorized",
		})
	}

	otp := s.otps.NewOTP()
	return c.JSON(model.UserLoginResponse{
		Status: "Ok",
		OTP:    otp,
	})
}

func NewSocketController(ctx context.Context, t *translator.UTtrans, v *validator.Validate, env config.SocketEnv) *socketController {
	c := &socketController{
		env:           env,
		activeClients: make(ActiveClient),
		translator:    t,
		validator:     v,
		handlers:      make(map[string]EventHandler),

		otps: otp.NewRetentionMap(ctx, env.OtpRetention),
	}
	c.setupEventHandlers()
	return c
}
