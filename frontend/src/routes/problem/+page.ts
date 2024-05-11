import type { SearchProblemResult } from "$lib/response";
import { fetchSearchProblemResult } from "./query";

export async function load({ url }: { url: URL }): Promise<SearchProblemResult> {
  const res = await fetchSearchProblemResult(url.searchParams);
  return res;
}
