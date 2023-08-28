import type { SubmissionResult } from "$lib/search";
import { API_HOST } from "$env/static/private";

export const csr = false;

export async function load({ url }) {
  const response = await fetch(`${API_HOST}/api/search/submission?${url.searchParams.toString()}`);
  const result: SubmissionResult = await response.json();

  return result;
}
