import type { SearchProblemResult } from "$lib/response";
import { fetchSearchProblemResult } from "./query";

export async function load({ url, fetch }): Promise<SearchProblemResult> {
  const res = await fetchSearchProblemResult(url.searchParams, fetch);
  return res;
}
