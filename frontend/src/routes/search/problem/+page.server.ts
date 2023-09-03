import type { ProblemSearchResult } from "$lib/search";
import { API_HOST } from "$env/static/private";

export const csr = false;

export async function load({ url }: { url: URL }) {
  const response = await fetch(`${API_HOST}/api/search/problem?${url.searchParams.toString()}`);
  const result: ProblemSearchResult = await response.json();

  return result;
}
