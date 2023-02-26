export type SearchResponse = {
  stats: Stats;
  items: Items;
};

export type Stats = {
  time: number;
  message: string | null;
  total: number;
  offset: number;
  amount: number;
  facet: Map<string, Facet>;
};

export type Items = {
  docs: Item[];
};

export type Item = {
  problem_id: string;
  problem_title: string;
  problem_url: string;
  contest_id: string;
  contest_title: string;
  contest_url: string;
  difficulty: number;
  start_at: Date;
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
