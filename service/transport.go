package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	model "youtube_service/model"
	db "youtube_service/repository"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
)

var errBadRoute, errInvalidRequest = errors.New("bad route"), errors.New("invalid request type")

// creating handlers for all the endpoints
func MakeHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}
	viewVideoHandler := kithttp.NewServer(
		MakeViewVideoEndpoint(s),
		decodeViewVideoRequest,
		encodeResponse,
		opts...,
	)
	GetViewsHandler := kithttp.NewServer(
		MakeGetViewsEndpoint(s),
		decodeGetViewsRequest,
		encodeResponse,
		opts...,
	)
	GetTopNVideosHandler := kithttp.NewServer(
		MakeGetTopNVideosEndpoint(s),
		decodeGetNvideosRequest,
		encodeResponse,
		opts...,
	)

	makeGetTopNVideosTodayHandler := kithttp.NewServer(
		MakeGetTopNVideosTodayEndpoint(s),
		decodeGetNvideosRequest,
		encodeResponse,
		opts...,
	)

	makePostVideoHandler := kithttp.NewServer(
		MakePostVideoEndpoint(s),
		decodePostVideoRequest,
		encodeResponse,
		opts...,
	)

	R := mux.NewRouter()
	R.Handle("/viewVideo", viewVideoHandler).Methods("GET")
	R.Handle("/getViews", GetViewsHandler).Methods("GET")
	R.Handle("/getTopNvideos", GetTopNVideosHandler).Methods("GET")
	R.Handle("/getTopNvideosToday", makeGetTopNVideosTodayHandler).Methods("GET")
	R.Handle("/postVideo", makePostVideoHandler).Methods("POST")

	return R

}

//server related

func decodeViewVideoRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	videoName := r.URL.Query().Get("videoName")
	if videoName == "" {
		return nil, errBadRoute
	}
	return viewVideoRequest{videoName: videoName}, nil
}

func decodeGetViewsRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	videoName := r.URL.Query().Get("videoName")
	if videoName == "" {
		return nil, errBadRoute
	}
	return getViewsRequest{videoName: videoName}, nil
}

func decodeGetNvideosRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	limit := r.URL.Query().Get("limit")
	num, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	return getTopNvideosRequest{limit: num}, nil
}

func decodePostVideoRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var body struct {
		VideoName string `json:"videoName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return postVideoRequest{videoName: body.VideoName}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// client related
func _Encode_viewVideo_Request(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/viewVideo"
	request1, ok := request.(viewVideoRequest)

	if ok {
		queryMap := req.URL.Query()
		queryMap.Add("videoName", request1.videoName)
		req.URL.RawQuery = queryMap.Encode()
		return nil
	} else {
		return errInvalidRequest
	}
}

func _Encode_GetTopNVideosEndpoint_Request(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/getTopNvideos"
	request1, ok := request.(getTopNvideosRequest)
	if ok {
		queryMap := req.URL.Query()
		queryMap.Add("limit", strconv.Itoa(request1.limit))
		req.URL.RawQuery = queryMap.Encode()
		return nil
	}
	return encodeRequest(ctx, req, request)
}

func _Encode_GetTopNVideosTodayEndpoint_Request(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/getTopNvideosToday"
	request1, ok := request.(getTopNVideosTodayRequest)
	if ok {
		queryMap := req.URL.Query()
		queryMap.Add("limit", strconv.Itoa(request1.limit))
		req.URL.RawQuery = queryMap.Encode()
		return nil
	}

	return errInvalidRequest
}

func _Encode_GetViewsEndpoint_Request(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/getViews"
	request1, ok := request.(getViewsRequest)
	if ok {
		queryMap := req.URL.Query()
		queryMap.Add("videoName", request1.videoName)
		req.URL.RawQuery = queryMap.Encode()
		return nil
	}
	return errInvalidRequest
}

func _Encode_PostVideoEndpoint_Request(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/postVideo"
	request1, ok := request.(postVideoRequest)
	if ok {
		queryMap := req.URL.Query()
		queryMap.Add("videoName", request1.videoName)
		req.URL.RawQuery = queryMap.Encode()
		return nil
	}
	return errInvalidRequest
}

// client's decoding functions
func _Decode_viewVideo_Response(ctx context.Context, resp *http.Response) (interface{}, error) {
	type viewVideoResponse1 struct {
		Response string `json:"reponse,omitempty"`
		Err      string `json:"error,omitempty"`
	}
	var response1 viewVideoResponse1
	err := json.NewDecoder(resp.Body).Decode(&response1)
	if err != nil {
		err = errors.New(response1.Err)
	}

	return viewVideoResponse{Response: response1.Response, Err: err}, err
}

func _Decode_GetTopNVideosEndpoint_Response(ctx context.Context, resp *http.Response) (interface{}, error) {
	type getNvideosResponse1 struct {
		TopVideos []model.ResultRedis `json:"TopVideos,omitempty"`
		Err       string              `json:"error,omitempty"`
	}
	var response1 getNvideosResponse1

	err := json.NewDecoder(resp.Body).Decode(&response1)
	if err != nil {
		err = errors.New(response1.Err)
	}

	return getTopNvideosResponse{TopVideos: response1.TopVideos, Err: err}, err
}

func _Decode_GetTopNVideosTodayEndpoint_Response(ctx context.Context, resp *http.Response) (interface{}, error) {

	var response getTopNVideosTodayResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func _Decode_GetViewsEndpoint_Response(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response getViewsResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func _Decode_PostVideoEndpoint_Response(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response postVideoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case db.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
