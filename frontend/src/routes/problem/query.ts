import { browser } from "$app/environment";
import { env } from "$env/dynamic/public";
import type { RecommendProblemParameter } from "$lib/request";
import { booleanFromQueryString, intoURLSearchParams, nullableBooleanFromQueryString, numberFromQueryString, type SearchProblemParameter } from "$lib/request";
import type { RecommendProblemResult, SearchProblemResult } from "$lib/response";
import { selections } from "./sort";

export async function fetchSearchProblemResult(params: URLSearchParams, fetch: (input: URL | RequestInfo, init?: RequestInit | undefined) => Promise<Response>): Promise<SearchProblemResult> {
  const p: SearchProblemParameter = {
    limit: 60,
    page: numberFromQueryString(params.get("p")) ?? 1,
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

  const host = browser ? String(env.PUBLIC_EXTERNAL_API_HOST) : String(env.PUBLIC_INTERNAL_API_HOST);
  const res = await fetch(`${host}/api/search/problem?${intoURLSearchParams(p).toString()}`);
  return await res.json();
}

export async function fetchRecommendProblemResult(
  params: RecommendProblemParameter,
  fetch: (input: URL | RequestInfo, init?: RequestInit | undefined) => Promise<Response>,
): Promise<RecommendProblemResult> {
  const host = browser ? String(env.PUBLIC_EXTERNAL_API_HOST) : String(env.PUBLIC_INTERNAL_API_HOST);
  const res = await fetch(`${host}/api/recommend/problem?${intoURLSearchParams(params).toString()}`);
  return await res.json();
}
