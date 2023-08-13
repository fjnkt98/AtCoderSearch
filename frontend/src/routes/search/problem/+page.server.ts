import type { ProblemSearchResult } from "$lib/search";

export async function load({ url }) {
  const response = await fetch(`http://localhost:8000/api/search/problem?${url.searchParams.toString()}`);
  const result: ProblemSearchResult = await response.json();

  return result;
}
