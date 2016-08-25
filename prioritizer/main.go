package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const schedulerName = "priority"

func main() {
	log.Println("Starting prioritizer...")

	doneChan := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	//go reconcileUnscheduledPods(2, doneChan, &wg)
	go prioritizePods(2, doneChan, &wg)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			log.Printf("Shutdown signal received, exiting...")
			close(doneChan)
			wg.Wait()
			os.Exit(0)
		}
	}
}
