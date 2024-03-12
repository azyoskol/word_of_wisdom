package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/azyoskol/word_of_wisdom/config"
	"github.com/azyoskol/word_of_wisdom/internal/log"
	"github.com/azyoskol/word_of_wisdom/internal/pow"
	"github.com/azyoskol/word_of_wisdom/internal/server"
	"golang.org/x/sync/errgroup"
)

func PrintAllEnvVariables() {
	allEnvs := os.Environ()

	for _, value := range allEnvs {
		log.Infow("Server", log.M{"env": value})
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatalw("Error loading configuration", log.M{"err": err})
	}

	powMiddleware, err := pow.NewPoW(cfg)
	if err != nil {
		log.Fatalw("Error creating PoW", log.M{"err": err})
	}

	tcpServer, err := server.NewPowServer(cfg, powMiddleware)
	if err != nil {
		log.Fatalw("Error creating TCP server", log.M{"err": err})
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return tcpServer.Startup(gCtx)
	})

	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Warnw("Stoping service because of SIGTERM", log.M{})
		cancel()
		return nil
	})

	err = g.Wait()
	if err != nil {
		log.Errorw("Stoping service because errors", log.M{"err": err})
	}
}
