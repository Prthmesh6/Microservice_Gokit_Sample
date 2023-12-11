package service

import (
	"context"
	"reflect"
	"testing"
	model "youtube_service/model"
	db "youtube_service/repository"
	mockDb "youtube_service/repository/mock"

	"github.com/golang/mock/gomock"
)

func Test_service_ViewVideo(t *testing.T) {
	// mocking repo
	ctr := gomock.NewController(t)

	// creating mock db
	newMockDB := mockDb.NewMockDatabase(ctr)

	newMockDB.EXPECT().IncreaseScore("video10", float64(1)).Times(1).Return(nil)

	type fields struct {
		database db.Database
	}
	type args struct {
		videoName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "basic",
			fields:  fields{database: newMockDB},
			args:    args{videoName: "video10"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				database: tt.fields.database,
			}
			if err := s.ViewVideo(context.Background(), tt.args.videoName); (err != nil) != tt.wantErr {
				t.Errorf("service.ViewVideo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_GetTopNVideos(t *testing.T) {
	ctr := gomock.NewController(t)
	newMockDB := mockDb.NewMockDatabase(ctr)
	newMockDB.EXPECT().GetSortedRecords(10, true).Times(1).Return([]model.ResultRedis{model.ResultRedis{VideoID: "video100", ViewCount: 104}}, nil)
	type fields struct {
		database db.Database
	}
	type args struct {
		n          int
		isLifeTime bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.ResultRedis
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "top 10 videos",
			fields:  fields{database: newMockDB},
			args:    args{n: 10, isLifeTime: true},
			want:    []model.ResultRedis{model.ResultRedis{VideoID: "video100", ViewCount: 104}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				database: tt.fields.database,
			}
			got, err := s.GetTopNVideos(context.Background(), tt.args.n, tt.args.isLifeTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetTopNVideos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetTopNVideos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetViews(t *testing.T) {
	ctr := gomock.NewController(t)
	newMockDB := mockDb.NewMockDatabase(ctr)
	newMockDB.EXPECT().GetScore("video500").Times(1).Return(float64(0), nil)

	type fields struct {
		database db.Database
	}
	type args struct {
		videoName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Get views for video 500",
			fields:  fields{database: newMockDB},
			args:    args{videoName: "video500"},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				database: tt.fields.database,
			}
			got, err := s.GetViews(context.Background(), tt.args.videoName)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetViews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("service.GetViews() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_PostVideo(t *testing.T) {
	ctr := gomock.NewController(t)
	newMockDB := mockDb.NewMockDatabase(ctr)
	newMockDB.EXPECT().Set("video500", float64(0)).Times(1)
	type fields struct {
		database db.Database
	}
	type args struct {
		videoName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "post a video video500",
			fields:  fields{database: newMockDB},
			args:    args{videoName: "video500"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				database: tt.fields.database,
			}
			if err := s.PostVideo(context.Background(), tt.args.videoName); (err != nil) != tt.wantErr {
				t.Errorf("service.PostVideo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
