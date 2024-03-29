export type SearchStats<F> = {
  time: number;
  total: number;
  index: number;
  pages: number;
  count: number;
  params: object;
  facet: F | null;
};

export type FacetPart = {
  label: string;
  count: number;
};

export type FilterRange = {
  from: string | null;
  to: string | null;
};

export type SearchResult<T, F> = {
  stats: SearchStats<F>;
  items: T[];
  message: string | null;
};

export type ProblemSearchResult = SearchResult<Problem, ProblemFacet>;

export type Problem = {
  problem_id: string;
  problem_title: string;
  problem_url: string;
  contest_id: string;
  contest_title: string;
  contest_url: string;
  difficulty: number | null;
  color: string | null;
  start_at: string;
  duration: number;
  rate_change: string;
  category: string;
};

export type ProblemFacet = {
  category: FacetPart[] | null;
  difficulty: FacetPart[] | null;
};

export type UserSearchResult = SearchResult<User, UserFacet>;

export type User = {
  user_name: string;
  rating: number;
  highest_rating: number;
  affiliation: string | null;
  birth_year: number | null;
  country: string | null;
  crown: string | null;
  join_count: number;
  rank: number;
  active_rank: number | null;
  wins: number;
  color: string;
  user_url: string;
};

export type UserFacet = {
  rating: FacetPart[] | null;
  birth_year: FacetPart[] | null;
  join_count: FacetPart[] | null;
  country: FacetPart[] | null;
};

export type SubmissionResult = SearchResult<Submission, SubmissionFacet>;

export type Submission = {
  submission_id: number;
  submitted_at: string;
  submission_url: string;
  problem_id: string;
  problem_title: string;
  contest_id: string;
  contest_title: string;
  category: string;
  difficulty: number;
  color: string;
  user_id: string;
  language: string;
  point: number;
  length: number;
  result: string;
  execution_time: number | null;
};

export type SubmissionFacet = {
  problem_id: FacetPart[] | null;
  user_id: FacetPart[] | null;
  language: FacetPart[] | null;
  result: FacetPart[] | null;
  length: FacetPart[] | null;
  execution_time: FacetPart[] | null;
};

export type RecommendResult = SearchResult<Recommend, object>;

export type Recommend = {
  problem_id: string;
  problem_title: string;
  problem_url: string;
  contest_id: string;
  contest_title: string;
  contest_url: string;
  difficulty: number | null;
  color: string | null;
  start_at: string;
  duration: number;
  rate_change: string;
  category: string;
  score: number;
};
