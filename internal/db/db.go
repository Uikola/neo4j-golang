package db

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
	"os"
)

func InitNeo4j(ctx context.Context) neo4j.DriverWithContext {
	driver, err := neo4j.NewDriverWithContext(
		os.Getenv("DB_URI"),
		neo4j.BasicAuth(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), ""))
	if err != nil {
		log.Error().Err(err).Msg("failed to create neo4j driver with context")
		os.Exit(1)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		log.Err(err).Msg("failed to verify neo4j db connectivity")
		os.Exit(1)
	}
	return driver
}
