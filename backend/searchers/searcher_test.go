package searchers

import (
	"context"
	pb "fjnkt98/atcodersearch/grpc/atcodersearch/v1"
	"testing"

	"github.com/meilisearch/meilisearch-go"
)

func TestSearchProblem(t *testing.T) {
	// _, url, key, stop, err := testutil.CreateEngineContainer()
	// t.Cleanup(func() { stop() })

	// if err != nil {
	// 	t.Fatal(err)
	// }

	client := meilisearch.New("http://localhost:7700", meilisearch.WithAPIKey("meili-master-key"))

	t.Run("search problem", func(t *testing.T) {
		ctx := context.Background()

		searcher := NewSearcher(client, nil)

		_, err := searcher.SearchProblem(ctx, &pb.SearchProblemRequest{})
		if err != nil {
			t.Error(err)
		}
	})
}
