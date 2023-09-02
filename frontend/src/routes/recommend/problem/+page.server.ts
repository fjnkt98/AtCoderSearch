import type { ProblemSearchResult, RecommendResult, UserSearchResult } from "$lib/search";
import { API_HOST } from "$env/static/private";
import type { Data } from "./data.d.ts";

export const csr = false;

async function fetchProblems() {
  const params = new URLSearchParams([
    ["limit", "10"],
    ["page", "1"],
    ["sort", "-start_at"],
  ]);
  const response = await fetch(`${API_HOST}/api/search/problem?${params.toString()}`);
  const result: ProblemSearchResult = await response.json();

  return result;
}

async function fetchUserRating(userId: string): Promise<number> {
  const params = new URLSearchParams([
    ["filter.user_id", userId],
    ["limit", "1"],
  ]);
  const response = await fetch(`${API_HOST}/api/search/user?${params.toString()}`);
  const result: UserSearchResult = await response.json();

  if (result.items.length != 1) {
    throw new Error("invalid user id");
  }
  return result.items[0].rating;
}

type UserWithRating = {
  userId: string;
  rating: number;
};

async function fetchRecommendByRating(url: URL, user: UserWithRating | null) {
  if (user == null) {
    return null;
  }

  const params = new URLSearchParams(url.searchParams);
  params.set("user_id", user.userId);
  params.set("rating", user.rating.toString());

  const response = await fetch(`${API_HOST}/api/recommend/problem?${params.toString()}`);
  const result: RecommendResult = await response.json();

  return result;
}

export async function load({ url }): Promise<Data> {
  let user: UserWithRating | null = null;
  const userId = url.searchParams.get("user_id");
  if (userId != null) {
    try {
      const rating = await fetchUserRating(userId);
      user = {
        userId: userId,
        rating: rating,
      };
    } catch (e) {
      user = null;
    }
  }

  const [recent, recByRating] = await Promise.all([fetchProblems(), fetchRecommendByRating(url, user)]);
  return {
    recent,
    recByRating,
  };
}
