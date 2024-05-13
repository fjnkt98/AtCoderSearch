import { env } from "$env/dynamic/public";
import type { RecommendProblemParameter } from "$lib/request";
import { booleanFromQueryString, intoURLSearchParams, nullableBooleanFromQueryString, numberFromQueryString, type SearchProblemParameter } from "$lib/request";
import type { RecommendProblemResult, SearchProblemResult } from "$lib/response";
import { selections } from "./sort";

export async function fetchSearchProblemResult(params: URLSearchParams, fetch: (input: URL | RequestInfo, init?: RequestInit | undefined) => Promise<Response>): Promise<SearchProblemResult> {
  const p: SearchProblemParameter = {
    limit: 20,
    page: 1,
    q: params.get("q"),
    sort: selections.get(params.get("s") ?? "2")?.values ?? ["-startAt"],
    facet: ["category", "difficulty"],
    category: params.getAll("category"),
    difficultyFrom: numberFromQueryString(params.get("difficultyFrom")),
    difficultyTo: numberFromQueryString(params.get("difficultyTo")),
    color: params.getAll("color"),
    userId: params.get("userId"),
    difficulty: numberFromQueryString(params.get("difficulty")),
    excludeSolved: booleanFromQueryString(params.get("excludeSolved")),
    experimental: nullableBooleanFromQueryString(params.get("experimental")),
    prioritizeRecent: booleanFromQueryString(params.get("prioritizeRecent")),
  };

  const res = await fetch(`${String(env.PUBLIC_API_HOST)}/api/search/problem?${intoURLSearchParams(p).toString()}`);
  return await res.json();
}

export async function fetchRecommendProblemResult(
  params: RecommendProblemParameter,
  fetch: (input: URL | RequestInfo, init?: RequestInit | undefined) => Promise<Response>,
): Promise<RecommendProblemResult> {
  const res = await fetch(`${String(env.PUBLIC_API_HOST)}/api/recommend/problem?${intoURLSearchParams(params).toString()}`);
  return await res.json();
}
