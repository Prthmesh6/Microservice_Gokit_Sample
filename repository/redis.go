package database

import (
	"context"
	"errors"
	"time"
	model "youtube_service/model"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var ErrUnknown = errors.New("not found")

type redisCache struct {
	client *redis.Client
	prefix string
}

func (*redisCache) getTodayKey(key string) string {
	return key + time.Now().Local().Format("2023-06-01")
}

func NewRedis(client *redis.Client, key string) *redisCache {
	return &redisCache{
		client: client,
		prefix: key,
	}
}

func (r *redisCache) CheckDBHealth() bool {
	// Ping the Redis server to check the connection
	pong, err := r.client.Ping(ctx).Result()
	if err != nil || pong != "PONG" {
		return false
	}
	return true
}

// Getting videos in sorted order of their view count
func (r *redisCache) GetSortedRecords(n int, isLifeTime bool) ([]model.ResultRedis, error) {
	var key string
	//Extracting the key as per requirement
	if isLifeTime {
		key = r.prefix
	} else {
		key = r.getTodayKey(r.prefix)
	}
	//Getting the videos sorted by view count from database
	redisResponse := r.client.ZRevRangeWithScores(ctx, key, 0, int64(n))
	responseArray, err := redisResponse.Result()
	if err != nil {
		return nil, err
	}
	//converting this redis response into business response type
	arrayResult := make([]model.ResultRedis, len(responseArray))
	for index := range arrayResult {
		arrayResult[index] = model.ResultRedis{
			VideoID:   responseArray[index].Member.(string),
			ViewCount: int(responseArray[index].Score),
		}
	}
	return arrayResult, nil
}

// Increasing the viewcount of the video by increasing it's score
func (r *redisCache) IncreaseScore(videoName string, increaseBy float64) (err error) {
	for _, key := range []string{r.prefix, r.getTodayKey(r.prefix)} {
		_, err := r.client.ZIncrBy(ctx, key, increaseBy, videoName).Result()

		if err != nil {
			return err
		}
	}
	return nil
}

// adding a new member score pair in database
func (r *redisCache) Set(member string, score float64) (err error) {
	_, err = r.client.ZAdd(ctx, r.prefix, &redis.Z{
		Score:  score,
		Member: member,
	}).Result()

	return err
}

// To get the views of a particular video
func (r *redisCache) GetScore(videoName string) (response float64, err error) {
	key := r.prefix
	response, err = r.client.ZScore(ctx, key, videoName).Result()
	if err != nil {
		return response, err
	}
	return response, err
}
