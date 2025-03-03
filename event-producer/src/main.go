package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"rest-server/api"
	"rest-server/logger"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	log := logger.NewLogger("main")
	log.Info("Starting application")

	myService := api.NewApiService()
	go myService.Start()

	<-c

	err := myService.Stop()
	if err != nil {
		log.Error("Error stopping service: %v", err.Error())
	}

	time.Sleep(1 * time.Second)
	log.Info("Shut down application")
}
