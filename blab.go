package blab

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkar/blab/server"
)

var (
	// interrupt signals Run to shutdown.
	interrupt = make(chan os.Signal, 1)
)

// Run will launch a server and wait for a kill sig to quit
func Run(conf *Config) error {
	s, err := server.New(conf.Host, conf.Port, conf.LogDir)
	if err != nil {
		return err
	}
	go s.Start()
	defer func() {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
		defer cancel()
		ch := make(chan struct{})
		go func() {
			s.Close()
			ch <- struct{}{}
		}()
		select {
		case <-ch:
			return
		case <-ctx.Done():
			log.Println("ERRO: timeout closing rooms")
			return
		}
	}()

	signal.Notify(interrupt)
	for {
		select {
		case sig := <-interrupt:
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL:
				return nil
			}
		}
	}
}
