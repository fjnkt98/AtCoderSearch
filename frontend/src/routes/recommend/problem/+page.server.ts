import type { ProblemSearchResult, RecommendResult } from "$lib/search";
import { env } from "$env/dynamic/private";
import type { Data } from "./data.d.ts";

export const csr = false;

async function fetchProblems() {
  const params = new URLSearchParams([
    ["limit", "10"],
    ["page", "1"],
    ["sort", "-start_at"],
  ]);
  const response = await fetch(`${env.API_HOST}/api/search/problem?${params.toString()}`);
  const result: ProblemSearchResult = await response.json();

  return result;
}

async function fetchRecommendByRating(url: URL, user: string | null) {
  if (user == null) {
    return null;
  }

  const params = new URLSearchParams(url.searchParams);
  params.set("user_id", user);
  const response = await fetch(`${env.API_HOST}/api/recommend/problem?${params.toString()}`);
  const result: RecommendResult = await response.json();

  return result;
}

export async function load({ url, cookies }): Promise<Data> {
  let userId: string | null = url.searchParams.get("user_id");
  if (userId == null) {
    userId = cookies.get("user_id") ?? null;
  } else {
    cookies.set("user_id", userId, { path: "/recommend/problem" });
  }

  const [recent, recByRating] = await Promise.all([fetchProblems(), fetchRecommendByRating(url, userId)]);
  return {
    recent,
    recByRating,
  };
}
