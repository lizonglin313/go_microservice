package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"rpc/go_kit_gRPC/proto"
	"strings"
)

// StringEndpoints 也要实现定义的Service接口.
// TODO: 暂时没想明白为什么...
// 应该是这样的：
// client.go 中的 NewStringClient 返回一个ServiceInterface
// 正是 StringEndpoints
type StringEndpoints struct {
	StringEndpoint      endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func (e StringEndpoints) Concat(ctx context.Context, a, b string) (string, error) {
	// TODO: 这里缺了一个 &  NOTICE 这就是问题所在！
	resp, err := e.StringEndpoint(ctx, &proto.StringRequest{A: a, B: b})
	//resp, err := e.StringEndpoint(ctx, StringRequest{RequestType: "Concat", A: a, B: b})
	//if resp == nil {
	//	fmt.Println("Resp is nil!")
	//}
	// TODO: 这里为什么要使用指针???
	// 经过调试,看起来和这里是不是指针都没关系
	// 正如值传递,我们直到 这里得到的 Response 是 nil
	response, ok := resp.(*proto.StringResponse)
	if !ok {
		fmt.Println("Get response->proto.StringResponse err, in endpoint.go Concat!")
	}
	//if response == nil {
	//	fmt.Println("response is nil, in endpoint.go Concat!")
	//}
	return response.Ret, err
}

//func (e StringEndpoints) Concat(ctx context.Context, a, b string) (string, error) {
//	resp, err := e.StringEndpoint(ctx, StringRequest{"Concat", a, b})
//	//if resp == nil {
//	//	fmt.Println("Resp is nil!")
//	//}
//	response := resp.(proto.StringResponse)
//	return response.Ret, err
//}

func (e StringEndpoints) Diff(ctx context.Context, a, b string) (string, error) {
	resp, err := e.StringEndpoint(ctx, &proto.StringRequest{A: a, B: b})
	// TODO: 不用指针可以吗?
	response := resp.(*proto.StringResponse)
	return response.Ret, err
}

var (
	ErrInvalidRequestType = errors.New("RequestType has only two type: Concat, Diff")
)

// StringRequest define request struct
type StringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

// StringResponse define response struct
type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

func MakeStringEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// 这里的 request 已经是 transport 处理好的正确格式的 request 了 所以直接转换
		req := request.(StringRequest)

		var (
			res, a, b string
			opError   error
		)

		a = req.A
		b = req.B

		fmt.Println("A in MakeStringEndpoint is:", a)

		if strings.EqualFold(req.RequestType, "Concat") {
			res, _ = svc.Concat(ctx, a, b)
		} else if strings.EqualFold(req.RequestType, "Diff") {
			res, _ = svc.Diff(ctx, a, b)
		} else {
			// fmt.Println("Error type!")
			return nil, ErrInvalidRequestType
		}
		// 把结果转换为 transport 可以接收的类型
		// fmt.Println("Result:", res)
		return StringResponse{Result: res, Error: opError}, nil
	}
}

type HealthRequest struct {
}

type HealthResponse struct {
	Status bool `json:"status"`
}

func MakeHealthCheckEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// 健康检查直接给他返回个 true
		return HealthResponse{true}, nil
	}
}
