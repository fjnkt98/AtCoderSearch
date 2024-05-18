import type { SearchUserResult } from "$lib/response";
import { fetchSearchUserResult } from "./query";

export async function load({ url, fetch }): Promise<SearchUserResult> {
  const res = await fetchSearchUserResult(url.searchParams, fetch);
  return res;
}
