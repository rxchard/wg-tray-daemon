package cmd

import (
	"context"
	"github.com/rxchard/wg-tray-daemon/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var rootx context.Context

func Execute() {
	var cancel context.CancelFunc
	rootx, cancel = context.WithCancel(context.Background())
	defer cancel()

	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-channel
		cancel()
	}()

	if err := server.Execute(rootx); err != nil {
		log.Fatal(err)
	}
}
