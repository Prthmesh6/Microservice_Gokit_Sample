package service

import (
	"context"
	"errors"
	model "youtube_service/model"
	db "youtube_service/repository"
)

var ErrInvalidArgument = errors.New("invalid argument")

type service struct {
	database db.Database
}

type Service interface {

	//viewVideo function is for viewing the particular video, it takes video name and increases view count by 1
	ViewVideo(context.Context, string) (err error)

	//get top N videos returns an array with top N videos with maximum views. it takes n integer and isLifeTime,
	//isLifeTime is for checking whether the request is for top N videos on Today or Top N videos for Lifetime
	//the returned array contains videoID and views
	GetTopNVideos(ctx context.Context, n int, isLifeTime bool) ([]model.ResultRedis, error)

	//getting the views for a particular video, this will return the total views any video have
	GetViews(ctx context.Context, videoName string) (int, error)

	//To add a new video PostVideo method will be used, it takes videoName and keeps the initial view count as zero
	PostVideo(ctx context.Context, videoName string) error
}

// create a new service by injecting a DB client
func NewService(database db.Database) Service {
	return &service{
		database: database,
	}
}

func (s *service) ViewVideo(ctx context.Context, videoName string) (err error) {
	if videoName == "" {
		return ErrInvalidArgument
	}
	s.increaseViewCount(ctx, videoName, 1)
	return err
}

func (s *service) GetTopNVideos(ctx context.Context, n int, isLifeTime bool) ([]model.ResultRedis, error) {
	if n == 0 {
		return nil, ErrInvalidArgument
	}
	arrayResult, err := s.database.GetSortedRecords(n, isLifeTime)
	if err != nil {
		return nil, err
	}
	return arrayResult, err
}

func (s *service) GetViews(ctx context.Context, videoName string) (int, error) {
	if videoName == "" {
		return -1, ErrInvalidArgument
	}
	views, err := s.database.GetScore(videoName)
	if err != nil {
		return int(views), err
	}
	return int(views), nil
}

func (s *service) PostVideo(ctx context.Context, videoName string) error {
	if videoName == "" {
		return ErrInvalidArgument
	}
	err := s.database.Set(videoName, 0)
	return err
}

// a function to increase a view count for a particular video
func (s *service) increaseViewCount(ctx context.Context, videoName string, increaseBy float64) error {
	if videoName == "" {
		return ErrInvalidArgument
	}
	err := s.database.IncreaseScore(videoName, increaseBy)
	if err != nil {
		return err
	}
	return nil
}
