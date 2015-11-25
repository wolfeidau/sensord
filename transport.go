package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

func makeRecordEndpoint(svc SensorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(recordRequest)
		err := svc.Record(req.I, req.V, req.T)
		if err != nil {
			return recordResponse{err.Error()}, nil
		}
		return recordResponse{""}, nil
	}
}

func decodeRecordRequest(r *http.Request) (interface{}, error) {
	var request recordRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type recordRequest struct {
	I string `json:"i"`
	V string `json:"v"`
	T string `json:"t"`
}

type recordResponse struct {
	Err string `json:"err,omitempty"`
}
