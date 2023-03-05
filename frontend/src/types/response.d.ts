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
  facet: Map<string, Facet>;
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

export type Facet = {
  counts: FacetCount[];
  start: string | null;
  end: string | null;
  gap: string | null;
  before: string | null;
  after: string | null;
  between: string | null;
};

export type FacetCount = {
  key: string;
  count: number;
};
