package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gpiecyk/data-warehouse/internal/api"
	"github.com/gpiecyk/data-warehouse/internal/configs"
	"github.com/gpiecyk/data-warehouse/internal/platform/cache"
	"github.com/gpiecyk/data-warehouse/internal/platform/database"
	"github.com/gpiecyk/data-warehouse/internal/server/http"
	"github.com/gpiecyk/data-warehouse/internal/users"
)

var (
	configPath = flag.String("config", "config.yaml", "configuration file path")
)

func main() {
	flag.Parse()

	err := validateConfigPath(*configPath)
	if err != nil {
		flag.Usage()
		logErrorIfPresent(err)
	}

	config, err := configs.NewService(*configPath)
	logErrorIfPresent(err)

	dbConfig, err := config.DatabaseConfig()
	logErrorIfPresent(err)

	db, err := database.NewService(dbConfig)
	logErrorIfPresent(err)

	cacheConfig, err := config.CacheConfig()
	logErrorIfPresent(err)

	cacheClient, err := cache.NewService(cacheConfig)
	logErrorIfPresent(err)

	userService, err := users.NewService(db, cacheClient)
	logErrorIfPresent(err)

	api, err := api.NewService(userService)
	logErrorIfPresent(err)

	httpConfig, err := config.HTTP()
	logErrorIfPresent(err)

	graphQL, err := config.GraphQLConfig(userService)
	logErrorIfPresent(err)

	httpServer, err := http.NewService(httpConfig, api, graphQL)
	logErrorIfPresent(err)

	httpServer.Start()
}

func validateConfigPath(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func logErrorIfPresent(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
