package setup

import (
	"log"
	"net/http"
	"os"
	config "youtube_service/config"
	db "youtube_service/repository"
	service "youtube_service/service"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/consul/api"

	log1 "github.com/go-kit/kit/log"
)

func SetUp() (http.Handler, *config.Config) {
	config_consul := api.DefaultConfig()
	client_consul, err := api.NewClient(config_consul)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}
	key, kv := "youtube_service_configs", client_consul.KV()

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		log.Printf("Failed to get key-value pair: %v\n setting up default configs", err)
	}
	if pair == nil {
		log.Printf("Key '%s' not found in Consul", key)
	}
	//setting up configs from consul
	configs := config.SetConfigs(pair)

	// Create a new Redis client
	var rdb = redis.NewClient(&redis.Options{
		Addr:     configs.RedisURL,
		Password: "",
		DB:       0,
	})

	redis := db.NewRedis(rdb, configs.RedisKey)

	var logger log1.Logger
	logger = log1.NewLogfmtLogger(log1.NewSyncWriter(os.Stderr))
	logger = log1.With(logger, "ts", log1.DefaultTimestampUTC)

	//creating a new service and wrapping it with logging layer
	yt_service := service.NewService(redis)
	yt_service = service.NewLoggingService(log1.With(logger), yt_service)

	mux := http.NewServeMux()
	mux.Handle("/", service.MakeHandler(yt_service, logger))
	return mux, configs

}
