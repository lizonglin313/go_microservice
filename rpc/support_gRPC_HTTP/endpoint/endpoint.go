package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"log"
	"rpc/support_gRPC_HTTP/service"
)

type Endpoints struct {
	SumEndpoint         endpoint.Endpoint
	ConcatEndpoint      endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func (e Endpoints) HealthCheck(ctx context.Context, status bool) bool {
	req := HealthCheckRequest{}
	resp, err := e.HealthCheckEndpoint(ctx, req)
	if err != nil {
		log.Println("Health Check Error!")
		return false
	}
	return resp.(HealthCheckResponse).Status
}

func (e Endpoints) Sum(ctx context.Context, a, b int) (reply int) {
	req := SumRequest{
		A: a,
		B: b,
	}
	resp, err := e.SumEndpoint(ctx, req)
	if err != nil {
		return -1
	}
	return resp.(SumResponse).Reply
}

func (e Endpoints) Concat(ctx context.Context, a, b string) (reply string) {
	req := ConcatRequest{
		A: a,
		B: b,
	}
	resp, err := e.ConcatEndpoint(ctx, req)
	if err != nil {
		return err.Error()
	}
	return resp.(ConcatResponse).Reply
}

type SumRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

type SumResponse struct {
	Reply int `json:"reply"`
}

func MakeSumEndpoint(svc service.SumConcatService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		edpSumReq := request.(SumRequest)
		reply := svc.Sum(ctx, edpSumReq.A, edpSumReq.B)
		return SumResponse{Reply: reply}, nil
	}
}

type ConcatRequest struct {
	A string `json:"a"`
	B string `json:"b"`
}

type ConcatResponse struct {
	Reply string `json:"reply"`
}

func MakeConcatEndpoint(svc service.SumConcatService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		edpConcatReq := request.(ConcatRequest)
		reply := svc.Concat(ctx, edpConcatReq.A, edpConcatReq.B)
		return ConcatResponse{Reply: reply}, nil
	}
}

type HealthCheckRequest struct {
}

type HealthCheckResponse struct {
	Status bool `json:"status"`
}

func MakeHealthCheckEndpoint(svc service.SumConcatService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck(ctx, true)
		return HealthCheckResponse{Status: status}, nil
	}
}
