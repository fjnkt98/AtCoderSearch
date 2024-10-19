package searchers

import (
	"context"
	pb "fjnkt98/atcodersearch/grpc/atcodersearch/v1"
	"log/slog"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meilisearch/meilisearch-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FieldList(doc any) []string {
	ty := reflect.TypeOf(doc)
	if ty.Kind() != reflect.Pointer {
		return nil
	}

	ty = ty.Elem()
	if ty.Kind() != reflect.Struct {
		return nil
	}

	fl := make([]string, 0, ty.NumField())
	for i := 0; i < ty.NumField(); i++ {
		f := ty.Field(i)

		var name string
		if tag, ok := f.Tag.Lookup("mapstructure"); ok {
			if tag == "-" {
				continue
			}
			n, _, _ := strings.Cut(tag, ",")
			name = n
		} else {
			name = f.Name
		}
		fl = append(fl, name)
	}
	return fl
}

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

type Problem struct {
	ProblemID      string  `mapstructure:"problemId" bun:"problem_id"`
	ProblemTitle   string  `mapstructure:"problemTitle" bun:"problem_title"`
	ProblemURL     string  `mapstructure:"problemUrl" bun:"problem_url"`
	ContestID      string  `mapstructure:"contestId" bun:"contest_id"`
	ContestTitle   string  `mapstructure:"contestTitle" bun:"contest_title"`
	ContestURL     string  `mapstructure:"contestUrl" bun:"contest_url"`
	Difficulty     *int64  `mapstructure:"difficulty" bun:"difficulty"`
	Color          *string `mapstructure:"color" bun:"color"`
	StartAt        int64   `mapstructure:"startAt" bun:"start_at"`
	Duration       int64   `mapstructure:"duration" bun:"duration"`
	RateChange     string  `mapstructure:"rateChange" bun:"rate_change"`
	Category       string  `mapstructure:"category" bun:"category"`
	IsExperimental bool    `mapstructure:"isExperimental" bun:"is_experimental"`
}

func (p *Problem) Into() *pb.Problem {
	return &pb.Problem{
		ProblemId:      p.ProblemID,
		ProblemTitle:   p.ProblemTitle,
		ProblemUrl:     p.ProblemURL,
		ContestId:      p.ContestID,
		ContestTitle:   p.ContestTitle,
		ContestUrl:     p.ContestURL,
		Difficulty:     p.Difficulty,
		Color:          p.Color,
		StartAt:        p.StartAt,
		Duration:       p.Duration,
		RateChange:     p.RateChange,
		Category:       p.Category,
		IsExperimental: p.IsExperimental,
	}
}

func (s *Searcher) SearchProblem(ctx context.Context, req *pb.SearchProblemRequest) (*pb.SearchProblemResponse, error) {
	panic("")
}

func (s *Searcher) SearchProblemByKeyword(ctx context.Context, req *pb.SearchProblemByKeywordRequest) (*pb.SearchProblemByKeywordResponse, error) {
	index := s.client.Index("problems")

	limit := req.GetLimit()
	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	q := &meilisearch.SearchRequest{
		AttributesToRetrieve: FieldList(new(Problem)),
		HitsPerPage:          int64(limit),
		Page:                 int64(page),
	}

	res, err := index.SearchWithContext(ctx, req.GetQ(), q)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "search failed", slog.Any("error", err))
		return nil, status.Errorf(codes.Unknown, "search: %s", err)
	}

	items := make([]*pb.Problem, len(res.Hits))
	for i, hit := range res.Hits {
		item, ok := hit.(map[string]any)
		if !ok {
			return nil, status.Errorf(codes.Unknown, "item conversion: item of res.Hits isn't an map[string]any")
		}

		var problem Problem
		if err := mapstructure.Decode(item, &problem); err != nil {
			return nil, status.Errorf(codes.Unknown, "item conversion: %s", err)
		}

		items[i] = problem.Into()
	}

	return &pb.SearchProblemByKeywordResponse{
		Time:  res.ProcessingTimeMs,
		Total: res.TotalHits,
		Index: res.Page,
		Pages: res.TotalPages,
		Items: items,
		Facet: nil,
	}, nil
}

func (s *Searcher) SearchUser(ctx context.Context, req *pb.SearchUserRequest) (*pb.SearchUserResponse, error) {
	panic("")
}

func (s *Searcher) SearchSubmission(ctx context.Context, req *pb.SearchSubmissionRequest) (*pb.SearchSubmissionResponse, error) {
	panic("")
}

func (s *Searcher) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.GetCategoryResponse, error) {
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
