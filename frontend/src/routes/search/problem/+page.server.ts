import type { ProblemSearchResult } from "$lib/search";
import { env } from "$env/dynamic/private";

export const csr = false;

export async function load({ url }: { url: URL }) {
  const response = await fetch(`${env.API_HOST}/api/search/problem?${url.searchParams.toString()}`);
  const result: ProblemSearchResult = await response.json();

  return result;
}
