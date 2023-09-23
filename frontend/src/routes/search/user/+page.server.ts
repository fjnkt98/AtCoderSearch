import type { UserSearchResult } from "$lib/search";
import { env } from "$env/dynamic/private";

export const csr = false;

export async function load({ url }) {
  const response = await fetch(`${env.API_HOST}/api/search/user?${url.searchParams.toString()}`);
  const result: UserSearchResult = await response.json();

  return result;
}
