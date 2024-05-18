import { env } from "$env/dynamic/public";
import { intoURLSearchParams, numberFromQueryString, type SearchUserParameter } from "$lib/request";
import type { SearchUserResult } from "$lib/response";
import { selections } from "./sort";

export async function fetchSearchUserResult(params: URLSearchParams, fetch: (input: URL | RequestInfo, init?: RequestInit | undefined) => Promise<Response>): Promise<SearchUserResult> {
  const p: SearchUserParameter = {
    limit: 60,
    page: numberFromQueryString(params.get("p")) ?? 1,
    q: params.get("q"),
    sort: selections.get(params.get("s") ?? "3")?.values ?? ["-rating"],
    facet: ["country", "rating", "birthYear", "joinCount"],
    userId: params.getAll("userId"),
    ratingFrom: numberFromQueryString(params.get("ratingFrom")),
    ratingTo: numberFromQueryString(params.get("ratingTo")),
    birthYearFrom: numberFromQueryString(params.get("birthYearFrom")),
    birthYearTo: numberFromQueryString(params.get("birthYearTo")),
    joinCountFrom: numberFromQueryString(params.get("joinCountFrom")),
    joinCountTo: numberFromQueryString(params.get("joinCountTo")),
    country: params.getAll("country"),
    color: params.getAll("color"),
  };

  const res = await fetch(`${String(env.PUBLIC_API_HOST)}/api/search/user?${intoURLSearchParams(p).toString()}`);
  return await res.json();
}
