package database

import (
	model "youtube_service/model"
)

//go:generate mockgen -source=db.go -destination mock/mock.go
type Database interface {
	Set(member string, score float64) error
	CheckDBHealth() bool
	GetScore(member string) (response float64, err error)
	GetSortedRecords(n int, ifLifeTime bool) ([]model.ResultRedis, error)
	IncreaseScore(videoName string, increaseBy float64) (err error)
}
