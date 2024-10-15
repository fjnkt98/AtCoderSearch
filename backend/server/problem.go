package server

import (
	"context"
	pb "fjnkt98/atcodersearch/grpc"
)

type ProblemSearcher struct {
	pb.UnimplementedProblemServiceServer
}

func NewProblemSearcher() *ProblemSearcher {
	return &ProblemSearcher{}
}

func (s *ProblemSearcher) SearchByKeyword(ctx context.Context, req *pb.SearchProblemByKeywordRequest) (*pb.SearchProblemResult, error) {
	return &pb.SearchProblemResult{
		Time:  0,
		Total: 0,
		Index: 0,
		Count: 0,
		Pages: 0,
		Items: nil,
		Facet: nil,
	}, nil
}
