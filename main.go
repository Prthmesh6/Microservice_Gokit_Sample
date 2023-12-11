package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"youtube_service/config"
	db "youtube_service/repository"
	service "youtube_service/service"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/consul/api"

	log1 "github.com/go-kit/kit/log"
)

const (
	defaultPort              = "8080"
	defaultRoutingServiceURL = "http://localhost:8080"
)

func main() {
	config_consul := api.DefaultConfig()
	client_consul, err := api.NewClient(config_consul)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}
	key, kv := "youtube_service_configs", client_consul.KV()

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		log.Printf("Failed to get key-value pair: %v", err)
		log.Println("-Setting default configs")

	}
	if pair == nil {
		log.Printf("Key '%s' not found in Consul\nsetting up default configs", key)
	}
	//setting up configs from consul
	config := config.SetConfigs(pair)

	// Create a new Redis client
	var rdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisURL,
		Password: "",
		DB:       0,
	})

	redis := db.NewRedis(rdb, config.RedisKey)

	var logger log1.Logger

	logger = log1.NewLogfmtLogger(log1.NewSyncWriter(os.Stderr))
	logger = log1.With(logger, "ts", log1.DefaultTimestampUTC)

	//creating a new service and wrapping it with logging layer
	yt_service := service.NewService(redis)
	yt_service = service.NewLoggingService(log1.With(logger), yt_service)

	mux := http.NewServeMux()
	mux.Handle("/", service.MakeHandler(yt_service, logger))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	//-----Graceful Shutdown------

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	go func() {
		log.Printf("Server started on %v \n", defaultRoutingServiceURL)
		log.Println(server.ListenAndServe())
	}()

	<-s
	shutDown(server)

}

func shutDown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Handle error while server shutdown")
	}
	log.Println("doing gracefull shutdown ")
}
