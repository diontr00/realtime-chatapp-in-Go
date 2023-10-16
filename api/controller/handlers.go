package controller

import (
	"realtime-chat/model"
)

func broadCast(e model.Event, c *socketClient) error {

	for client := range c.controller.activeClients {
		if client != c {
			client.egress <- e.Payload
		}
	}
	return nil

}
