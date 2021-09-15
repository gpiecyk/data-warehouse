package configs

import (
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/gpiecyk/data-warehouse/graph"
	"github.com/gpiecyk/data-warehouse/graph/generated"
	"github.com/gpiecyk/data-warehouse/internal/platform/cache"
	"github.com/gpiecyk/data-warehouse/internal/platform/database"
	"github.com/gpiecyk/data-warehouse/internal/server/http"
	"github.com/gpiecyk/data-warehouse/internal/users"
	"gopkg.in/yaml.v3"
)

type Configs struct {
	Server struct {
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		ReadTimeout  int    `yaml:"readTimeout"`
		WriteTimeout int    `yaml:"writeTimeout"`
	}
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
		TimeZone string `yaml:"TimeZone"`
	}
	Cache struct {
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		PoolSize     int    `yaml:"poolSize"`
		IdleTimeout  int    `yaml:"idleTimeout"`
		ReadTimeout  int    `yaml:"readTimeout"`
		WriteTimeout int    `yaml:"writeTimeout"`
		DialTimeout  int    `yaml:"dialTimeout"`
	}
}

func (configs *Configs) HTTP() (*http.Config, error) {
	return &http.Config{
		Host:         configs.Server.Host,
		Port:         configs.Server.Port,
		ReadTimeout:  time.Duration(configs.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(configs.Server.WriteTimeout) * time.Second,
	}, nil
}

func (configs *Configs) DatabaseConfig() (*database.Config, error) {
	return &database.Config{
		Host:     configs.Database.Host,
		Port:     configs.Database.Port,
		User:     configs.Database.User,
		Password: configs.Database.Password,
		Dbname:   configs.Database.Dbname,
		SSLMode:  configs.Database.SSLMode,
		TimeZone: configs.Database.TimeZone,
	}, nil
}

func (configs *Configs) CacheConfig() (*cache.Config, error) {
	return &cache.Config{
		Host:         configs.Cache.Host,
		Port:         configs.Cache.Port,
		PoolSize:     configs.Cache.PoolSize,
		IdleTimeout:  time.Duration(configs.Cache.IdleTimeout) * time.Second,
		ReadTimeout:  time.Duration(configs.Cache.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(configs.Cache.WriteTimeout) * time.Second,
		DialTimeout:  time.Duration(configs.Cache.DialTimeout) * time.Second,
	}, nil
}

func (configs *Configs) GraphQLConfig(service *users.UserService) (*handler.Server, error) {
	graphQL := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					UserService: service,
				},
			}))
	graphQL.Use(extension.FixedComplexityLimit(100))
	return graphQL, nil
}

func NewService(configPath string) (*Configs, error) {
	configs, err := getConfigsFromFile(configPath)
	return configs, err
}

func getConfigsFromFile(path string) (*Configs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	configs := &Configs{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(configs); err != nil {
		return nil, err
	}

	return configs, nil
}
