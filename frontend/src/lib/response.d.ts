export type ResultStats<F> = {
  time: number;
  total: number;
  index: number;
  pages: number;
  count: number;
  params: object | null;
  facet: F | null;
};

export type FacetCount = {
  label: string;
  count: number;
};

export type ResultResponse<T, F> = {
  stats: ResultStats<F>;
  items: T[];
  message: string | null;
};

export type SearchProblemFacet = {
  category?: FacetCount[];
  difficulty?: FacetCount[];
};

export type SearchProblemResult = ResultResponse<Problem, SearchProblemFacet>;
export type RecommendProblemResult = ResultResponse<Problem, SearchProblemFacet>;

export type Problem = {
  problemId: string;
  problemTitle: string;
  problemUrl: string;
  contestId: string;
  contestTitle: string;
  contestUrl: string;
  difficulty: number | null;
  color: string | null;
  startAt: string;
  duration: number;
  rateChange: string;
  category: string;
  score: number;
};

export type SearchUserResult = ResultResponse<User, SearchUserFacet>;

export type User = {
  userId: string;
  rating: number;
  highestRating: number;
  affiliation: string | null;
  birthYear: number | null;
  country: string | null;
  crown: string | null;
  joinCount: number;
  rank: number;
  activeRank: number | null;
  wins: number;
  color: string;
  userUrl: string;
};

export type SearchUserFacet = {
  country?: FacetCount[];
  rating?: FacetCount[];
  birthYear?: FacetCount[];
  joinCount?: FacetCount[];
};

export type SearchSubmissionResult = ResultResponse<Submission>;

export type Submission = {
  submissionId: number;
  submittedAt: string;
  submissionUrl: string;
  problemId: string;
  problemTitle: string;
  contestId: string;
  contestTitle: string;
  category: string;
  difficulty: number;
  color: string;
  userId: string;
  language: string;
  point: number;
  length: number;
  result: string;
  executionTime: number | null;
};
