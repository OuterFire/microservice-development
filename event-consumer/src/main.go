package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"event-consumer/consumer"
	"event-consumer/logger"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	log := logger.NewLogger("main")
	log.Info("Starting application")

	rd := consumer.NewConsumerService()
	rd.Start()

	<-c

	err := rd.Stop()
	if err != nil {
		log.Error("Error stopping consumer: %v", err)
	}
	time.Sleep(2 * time.Second)
	log.Info("Shut down application")
}
