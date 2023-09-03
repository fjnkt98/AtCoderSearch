import type { RecommendResult, ProblemSearchResult } from "$lib/search";

export type Data = {
  recent: ProblemSearchResult;
  recByRating: RecommendResult | null;
};
