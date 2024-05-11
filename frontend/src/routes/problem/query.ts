import { env } from "$env/dynamic/public";
import { intoURLSearchParams, numberFromQueryString, type SearchProblemParameter } from "$lib/request";
import type { SearchProblemResult } from "$lib/response";
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
    excludeSolved: Boolean(params.get("excludeSolved")),
    experimental: Boolean(params.get("experimental")),
    prioritizeRecent: Boolean(params.get("prioritizeRecent")),
  };

  const res = await fetch(`${String(env.PUBLIC_API_HOST)}/api/search/problem?${intoURLSearchParams(p).toString()}`);
  return await res.json();
}
