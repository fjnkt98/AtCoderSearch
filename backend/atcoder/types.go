package atcoder

type SubmissionList struct {
	MaxPage     int
	Submissions []Submission
}

type Submission struct {
	ID            int64   `db:"id"`
	EpochSecond   int64   `db:"epoch_second"`
	ProblemID     string  `db:"problem_id"`
	ContestID     string  `db:"contest_id"`
	UserID        string  `db:"user_id"`
	Language      string  `db:"language"`
	Point         float64 `db:"point"`
	Length        int64   `db:"length"`
	Result        string  `db:"result"`
	ExecutionTime *int64  `db:"execution_time"`
}
