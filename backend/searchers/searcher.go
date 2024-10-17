package searchers

import (
	"context"
	pb "fjnkt98/atcodersearch/grpc/atcodersearch/v1"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meilisearch/meilisearch-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Searcher struct {
	pb.UnimplementedSearchServiceServer

	client meilisearch.ServiceManager
	pool   *pgxpool.Pool
}

func NewSearcher(client meilisearch.ServiceManager, pool *pgxpool.Pool) *Searcher {
	return &Searcher{
		client: client,
		pool:   pool,
	}
}

func (s *Searcher) SearchProblem(ctx context.Context, req *pb.SearchProblemRequest) (*pb.SearchProblemResponse, error) {
	index := s.client.Index("problems")

	limit := req.GetLimit()
	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	q := &meilisearch.SearchRequest{
		AttributesToRetrieve: []string{
			"problemId",
			"problemTitle",
			"problemUrl",
			"contestId",
			"contestTitle",
			"contestUrl",
			"difficulty",
			"color",
			"startAt",
			"duration",
			"rateChange",
			"category",
		},
		HitsPerPage: int64(limit),
		Page:        int64(page),
	}

	result, err := index.SearchWithContext(ctx, req.GetQ(), q)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "search failed", slog.Any("error", err))
		return nil, status.Errorf(codes.Unknown, "search failed")
	}

	fmt.Printf("%+v\n", result.Hits)

	return &pb.SearchProblemResponse{
		Time:  0,
		Total: 0,
		Index: 0,
		Pages: 0,
		Items: nil,
		Facet: nil,
	}, nil
}

func (s *Searcher) SearchUser(ctx context.Context, req *pb.SearchUserRequest) (*pb.SearchUserResponse, error) {
	panic("")
}

func (s *Searcher) SearchSubmission(ctx context.Context, req *pb.SearchSubmissionRequest) (*pb.SearchSubmissionResponse, error) {
	panic("")
}

func (s *Searcher) GetCategory(ctx context.Context, req *emptypb.Empty) (*pb.GetCategoryResponse, error) {
	panic("")
}

func (s *Searcher) GetContest(ctx context.Context, req *pb.GetContestRequest) (*pb.GetContestResponse, error) {
	panic("")
}

func (s *Searcher) GetLanguage(ctx context.Context, req *pb.GetLanguageRequest) (*pb.GetLanguageResponse, error) {
	panic("")
}

func (s *Searcher) GetProblem(ctx context.Context, req *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	panic("")
}
