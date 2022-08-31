package main

import (
	"fmt"
	"myself_framwork/routers"
	"myself_framwork/utils"
	"os"
	"os/signal"
	"sync"

	"syscall"
)

// main ..
func main() {

	var wg sync.WaitGroup
	var errChan = make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		listenAddress := utils.GetEnv("API_SERVICE", "0.0.0.0:4000")
		fmt.Println("Starting listen address: ", listenAddress)
		errChan <- routers.Server(listenAddress)
	}()
	wg.Wait()

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		fmt.Println("Got an interrupt, exiting...")
	case err := <-errChan:
		if err != nil {
			fmt.Println("Error while running api, exiting...: ", err)
		}
	}
}
