package service

import (
	"context"
	"net/url"
	"strings"

	model "youtube_service/model"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

//Endpoints are the funtions which take request and return the response, we have client side endpoints and server side endpoints
//client side endpoints will take arguments, encode them into an http request, call the server endpoint, get the response
//decode the response in the funtion return types and return this in the client
//server side endpoints will take the request, will decode it using the decode function, call the business logic using service interface
//get the response from service and encode this into http response and return to the particular request

type Endpoints struct {
	ViewVideoEndpoint          endpoint.Endpoint
	GetTopNVideosEndpoint      endpoint.Endpoint
	GetTopNVideosTodayEndpoint endpoint.Endpoint
	GetViewsEndpoint           endpoint.Endpoint
	PostVideoEndpoint          endpoint.Endpoint
}

//kept for future use

// func makeServerEndpoints(s Service) *Endpoints {
// 	return &Endpoints{
// 		ViewVideoEndpoint:          MakeViewVideoEndpoint(s),
// 		GetTopNVideosEndpoint:      MakeGetTopNVideosEndpoint(s),
// 		GetTopNVideosTodayEndpoint: MakeGetTopNVideosTodayEndpoint(s),
// 		GetViewsEndpoint:           MakeGetViewsEndpoint(s),
// 		PostVideoEndpoint:          MakePostVideoEndpoint(s),
// 	}
// }

func (e Endpoints) ViewVideo(ctx context.Context, videoName string) (err error) {
	req := viewVideoRequest{videoName: videoName}
	response, err := e.ViewVideoEndpoint(context.Background(), req)
	if err != nil {
		return err
	}
	resp := response.(viewVideoResponse)
	return resp.Err
}

// if user want for lifefime we will hit the endpoint for lifetime top videos
// else we will hit the endpoint for top videos on current day
func (e Endpoints) GetTopNVideos(ctx context.Context, n int, isLifeTime bool) ([]model.ResultRedis, error) {
	var response interface{}
	var err error
	if isLifeTime {
		req := getTopNvideosRequest{limit: n}
		response, err = e.GetTopNVideosEndpoint(context.Background(), req)
	} else {
		req := getTopNVideosTodayRequest{limit: n}
		response, err = e.GetTopNVideosTodayEndpoint(context.Background(), req)
	}
	if err != nil {
		return nil, err
	}
	resp := response.(getTopNvideosResponse)
	return resp.TopVideos, resp.Err
}

func (e Endpoints) GetViews(ctx context.Context, videoName string) (int, error) {
	req := getViewsRequest{videoName: videoName}
	response, err := e.GetViewsEndpoint(context.Background(), req)
	if err != nil {
		return 0, err
	}
	reponse1 := response.(getViewsResponse)
	return reponse1.Views, reponse1.Err
}

func (e Endpoints) PostVideo(ctx context.Context, videoName string) error {
	req := postVideoRequest{videoName: videoName}
	_, err := e.PostVideoEndpoint(context.Background(), req)
	if err != nil {
		return err
	}
	resp := postVideoResponse{Err: err}
	return resp.Err

}

// httptransport.NewClient().endpoint() will create an endpoint by taking encoder decoder functions, target URL, request type and options
// and will return an usable client endpoint which calls the remote HTTP endpoint
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}

	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	return Endpoints{
		ViewVideoEndpoint:          httptransport.NewClient("GET", tgt, _Encode_viewVideo_Request, _Decode_viewVideo_Response, options...).Endpoint(),
		GetTopNVideosEndpoint:      httptransport.NewClient("GET", tgt, _Encode_GetTopNVideosEndpoint_Request, _Decode_GetTopNVideosEndpoint_Response, options...).Endpoint(),
		GetTopNVideosTodayEndpoint: httptransport.NewClient("GET", tgt, _Encode_GetTopNVideosTodayEndpoint_Request, _Decode_GetTopNVideosTodayEndpoint_Response, options...).Endpoint(),
		GetViewsEndpoint:           httptransport.NewClient("GET", tgt, _Encode_GetViewsEndpoint_Request, _Decode_GetViewsEndpoint_Response, options...).Endpoint(),
		PostVideoEndpoint:          httptransport.NewClient("POST", tgt, _Encode_PostVideoEndpoint_Request, _Decode_PostVideoEndpoint_Response, options...).Endpoint(),
	}, nil
}

//creating endpoints for all the methods in service

type viewVideoResponse struct {
	Response string `json:"reponse,omitempty"`
	Err      error  `json:"error,omitempty"`
}

type viewVideoRequest struct {
	videoName string
}

func (r viewVideoResponse) error() error { return r.Err }

func MakeViewVideoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewVideoRequest)
		err = s.ViewVideo(ctx, req.videoName)
		return viewVideoResponse{Response: "success", Err: err}, nil
	}
}

type getTopNvideosRequest struct {
	limit int
}

type getTopNvideosResponse struct {
	TopVideos []model.ResultRedis
	Err       error
}

func (r getTopNvideosResponse) error() error { return r.Err }

func MakeGetTopNVideosEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getTopNvideosRequest)
		isLifeTime := true
		response, err = s.GetTopNVideos(ctx, req.limit, isLifeTime)
		response1, _ := response.([]model.ResultRedis)
		return getTopNvideosResponse{TopVideos: response1, Err: err}, nil
	}
}

type getTopNVideosTodayRequest struct {
	limit int
}

type getTopNVideosTodayResponse struct {
	TopVideos []model.ResultRedis
	Err       error
}

func (r getTopNVideosTodayResponse) error() error { return r.Err }

func MakeGetTopNVideosTodayEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getTopNVideosTodayRequest)
		isLifeTime := false
		response, err = s.GetTopNVideos(ctx, req.limit, isLifeTime)
		response1, _ := response.([]model.ResultRedis)
		return getTopNVideosTodayResponse{TopVideos: response1, Err: err}, nil
	}
}

type getViewsRequest struct {
	videoName string
}

type getViewsResponse struct {
	Views int
	Err   error
}

func (r getViewsResponse) error() error { return r.Err }

func MakeGetViewsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getViewsRequest)
		response, err = s.GetViews(ctx, req.videoName)
		response1, _ := response.(int)
		return getViewsResponse{Views: response1, Err: err}, nil
	}
}

type postVideoResponse struct {
	Err error
}

type postVideoRequest struct {
	videoName string
}

func (r postVideoResponse) error() error { return r.Err }

func MakePostVideoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postVideoRequest)
		err = s.PostVideo(ctx, req.videoName)
		return postVideoResponse{Err: err}, nil
	}
}
