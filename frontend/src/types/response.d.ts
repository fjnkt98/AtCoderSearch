export type SearchResponse = {
  stats: Stats;
  items: Item[];
  message: string | null;
};

export type Stats = {
  time: number;
  total: number;
  index: number;
  pages: number;
  count: number;
  facet: FacetResults;
};

export type Item = {
  problem_id: string;
  problem_title: string;
  problem_url: string;
  contest_id: string;
  contest_title: string;
  contest_url: string;
  difficulty: number;
  start_at: string;
  duration: number;
  rate_change: string;
  category: string;
};

export type FacetResults = {
  category: FieldFacetResult;
  difficulty: RangeFacetResult;
};

export type FieldFacetResult = {
  counts: FieldFacetCount[];
};

export type FieldFacetCount = {
  key: string;
  count: number;
};

export type RangeFacetResult = {
  counts: RangeFacetCount[];
  start: string;
  end: string;
  gap: string;
  before: string | null;
  after: string | null;
  between: string | null;
};

export type RangeFacetCount = {
  begin: string;
  end: string;
  count: number;
};
