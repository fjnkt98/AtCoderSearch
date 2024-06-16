import type { SearchSubmissionResult } from "$lib/response";
import { fetchSearchSubmissionResult } from "./query";

export async function load({ url, fetch }): Promise<SearchSubmissionResult> {
  const res = await fetchSearchSubmissionResult(url.searchParams, fetch);
  return res;
}
