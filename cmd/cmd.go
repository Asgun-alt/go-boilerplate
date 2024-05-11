package cmd

import (
	"go-boilerplate/delivery/container"
	"go-boilerplate/delivery/server"
)

func Execute() {
	cont := container.New()

	server.New(cont)
}
