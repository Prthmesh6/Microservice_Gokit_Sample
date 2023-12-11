package service

import (
	"context"
	"time"
	model "youtube_service/model"

	log1 "github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log1.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log1.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) ViewVideo(ctx context.Context, videoName string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ViewVideo",
			"videoName", videoName,
			"err", err,
		)
	}(time.Now())
	return s.Service.ViewVideo(ctx, videoName)
}

func (s *loggingService) GetTopNVideos(ctx context.Context, n int, isLifeTime bool) (arraylist []model.ResultRedis, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "GetTopNVideos",
			"N", n,
			"isLifeTime", isLifeTime,
			"error", err,
		)
	}(time.Now())
	return s.Service.GetTopNVideos(ctx, n, isLifeTime)
}

func (s *loggingService) GetViews(ctx context.Context, videoName string) (int, error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "GetViews",
			"videoName", videoName,
		)
	}(time.Now())
	return s.Service.GetViews(ctx, videoName)
}

func (s *loggingService) PostVideo(ctx context.Context, videoName string) error {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "PostVideo",
			"videoName", videoName,
		)
	}(time.Now())
	return s.Service.PostVideo(ctx, videoName)
}
