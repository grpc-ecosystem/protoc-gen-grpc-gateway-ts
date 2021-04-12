package main

import (
	"context"
	"time"
)

type RealCounterService struct{}

func (r *RealCounterService) ExternalMessage(ctx context.Context, request *ExternalRequest) (*ExternalResponse, error) {
	return &ExternalResponse{
		Result: request.Content + "!!",
	}, nil
}

func (r *RealCounterService) Increment(c context.Context, req *UnaryRequest) (*UnaryResponse, error) {
	return &UnaryResponse{
		Result: req.Counter + 1,
	}, nil
}

func (r *RealCounterService) StreamingIncrements(req *StreamingRequest, service CounterService_StreamingIncrementsServer) error {
	times := 5
	counter := req.Counter

	for i := 0; i < times; i++ {
		counter++
		err := service.Send(&StreamingResponse{
			Result: counter,
		})
		if err != nil {
			return err
		}

		time.Sleep(200 * time.Millisecond)
	}

	return nil
}

func (r *RealCounterService) HTTPGet(ctx context.Context, req *HttpGetRequest) (*HttpGetResponse, error) {
	return &HttpGetResponse{
		Result: req.NumToIncrease + 1,
	}, nil
}

func (r *RealCounterService) HTTPPostWithNestedBodyPath(ctx context.Context, in *HttpPostRequest) (*HttpPostResponse, error) {
	return &HttpPostResponse{
		PostResult: in.A + in.Req.B,
	}, nil
}

func (r *RealCounterService) HTTPPostWithStarBodyPath(ctx context.Context, in *HttpPostRequest) (*HttpPostResponse, error) {
	return &HttpPostResponse{
		PostResult: in.A + in.Req.B + in.C,
	}, nil
}

func (r *RealCounterService) HTTPPatch(ctx context.Context, in *HttpPatchRequest) (*HttpPatchResponse, error) {
	return &HttpPatchResponse{
		Patch: in.A + in.C,
	}, nil
}
