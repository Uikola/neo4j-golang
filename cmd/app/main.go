package main

import (
	"context"
	"errors"
	"github.com/Uikola/neo4j-golang/internal/db"
	"github.com/Uikola/neo4j-golang/internal/db/repository/neo4j"
	server "github.com/Uikola/neo4j-golang/internal/server/http"
	"github.com/Uikola/neo4j-golang/internal/server/http/host"
	"github.com/Uikola/neo4j-golang/pkg/zlog"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

const (
	DEBUGLEVEL = 0
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}

	log.Logger = zlog.Default(true, "dev", DEBUGLEVEL)

	neo4jDriver := db.InitNeo4j(ctx)
	defer neo4jDriver.Close(ctx)

	hostRepository := neo4j.NewRepository(neo4jDriver)
	hostHandler := host.NewHandler(hostRepository)

	srv := server.NewServer(hostHandler)

	httpServer := &http.Server{
		Addr:    ":8001",
		Handler: srv,
	}

	go func() {
		log.Info().Msg("Starting server...")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Error().Msg(err.Error())
			os.Exit(1)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Error().Msg(err.Error())
			os.Exit(1)
		}
	}()
	wg.Wait()
}
