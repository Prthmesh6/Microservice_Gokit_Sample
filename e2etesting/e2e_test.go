package e2etesting

import (
	"context"
	"net/http/httptest"
	"testing"
	service "youtube_service/service"
	"youtube_service/setup"
)

func Test_service_PostVideo(t *testing.T) {
	mux, _ := setup.SetUp()
	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	endpoints, err := service.MakeClientEndpoints(testServer.URL)
	if err != nil {
		return
	}

	err = endpoints.PostVideo(context.Background(), "video10")
	if err != nil {
		t.Errorf("got %v while posting video", err)
	} else {
		t.Logf("testcase passed video posted successfully")
	}
}

func Test_service_ViewVideo(t *testing.T) {

	mux, _ := setup.SetUp()
	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	endpoints, err := service.MakeClientEndpoints(testServer.URL)
	if err != nil {
		return
	}
	err = endpoints.ViewVideo(context.Background(), "video10")
	if err != nil {
		t.Errorf("Got error while viewing the video %+v", err)
	} else {
		t.Logf("Testcase passed successfully")
	}
}

func Test_service_GetTopNVideos(t *testing.T) {
	mux, _ := setup.SetUp()
	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	endpoints, err := service.MakeClientEndpoints(testServer.URL)
	if err != nil {
		return
	}

	response, err := endpoints.GetTopNVideos(context.Background(), 10, true)
	if err != nil {
		t.Errorf("got %v while getting top N videos", err)
	} else if len(response) < 1 {
		t.Errorf("expected top 10 videos got %v", len(response))
	} else {
		t.Logf("Testcase passed successfully got top 10 videos %v", response)
	}
}

func Test_service_GetViews(t *testing.T) {
	mux, _ := setup.SetUp()
	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	endpoints, err := service.MakeClientEndpoints(testServer.URL)
	if err != nil {
		return
	}

	response, err := endpoints.GetViews(context.Background(), "video10")
	if err != nil {
		t.Errorf("got %v while getting views for a video10", err)
	} else if response < 1 {
		t.Errorf("expected response 1 got %v", response)
	} else {
		t.Logf("testcase passed expected 66 got %v", response)
	}
}
