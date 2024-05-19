import { env } from "$env/dynamic/public";
import { intoURLSearchParams, numberFromQueryString, type SearchSubmissionParameter } from "$lib/request";
import type { SearchSubmissionResult } from "$lib/response";
import { selections } from "./sort";

export async function fetchSearchSubmissionResult(params: URLSearchParams, fetch: (input: URL | RequestInfo, init?: RequestInit | undefined) => Promise<Response>): Promise<SearchSubmissionResult> {
  const p: SearchSubmissionParameter = {
    limit: 100,
    page: numberFromQueryString(params.get("p")) ?? 1,
    sort: selections.get(params.get("s") ?? "3")?.values ?? ["-submittedAt"],
    epochSecondFrom: numberFromQueryString(params.get("epochSecondFrom")),
    epochSecondTo: numberFromQueryString(params.get("epochSecondTo")),
    problemId: params.getAll("problemId"),
    contestId: params.getAll("contestId"),
    category: params.getAll("category"),
    userId: params.getAll("userId"),
    language: params.getAll("language"),
    languageGroup: params.getAll("languageGroup"),
    pointFrom: numberFromQueryString(params.get("pointFrom")),
    pointTo: numberFromQueryString(params.get("pointTo")),
    lengthFrom: numberFromQueryString(params.get("lengthFrom")),
    lengthTo: numberFromQueryString(params.get("lengthTo")),
    result: params.getAll("result"),
    executionTimeFrom: numberFromQueryString(params.get("executionTimeFrom")),
    executionTimeTo: numberFromQueryString(params.get("executionTimeTo")),
  };

  const res = await fetch(`${String(env.PUBLIC_API_HOST)}/api/search/submission?${intoURLSearchParams(p).toString()}`);
  return await res.json();
}
