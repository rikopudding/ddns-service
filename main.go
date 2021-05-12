package main

import (
	"ddns-service/config"
	"ddns-service/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	config.Init()

}

func main() {
	service.Start()
	signalHandler()
}

func signalHandler() {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Println("ddns-service exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
