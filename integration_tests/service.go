package main

import (
	"context"
	"time"
)

type RealCounterService struct{}

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
