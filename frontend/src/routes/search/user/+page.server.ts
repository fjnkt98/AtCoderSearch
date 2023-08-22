import type { UserSearchResult } from "$lib/search";
import { API_HOST } from "$env/static/private";

export const csr = false;

export async function load({ url }) {
  const response = await fetch(`${API_HOST}/api/search/user?${url.searchParams.toString()}`);
  const result: UserSearchResult = await response.json();

  return result;
}
