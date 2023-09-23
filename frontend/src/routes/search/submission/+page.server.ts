import type { SubmissionResult } from "$lib/search";
import { env } from "$env/dynamic/private";
import type { Data } from "./data.js";

export const csr = false;

async function fetchResult(url: URL): Promise<SubmissionResult> {
  const response = await fetch(`${env.API_HOST}/api/search/submission?${url.searchParams.toString()}`);
  const result = await response.json();

  return result;
}

async function fetchCategories(): Promise<string[]> {
  const response = await fetch(`${env.API_HOST}/api/list/category`);
  const result = await response.json();

  return result;
}

async function fetchLanguages(): Promise<string[]> {
  const response = await fetch(`${env.API_HOST}/api/list/language`);
  const result = await response.json();

  return result;
}

async function fetchContests(category: string | null): Promise<string[]> {
  const url = category != null ? `${env.API_HOST}/api/list/contest?category=${category}` : `${env.API_HOST}/api/list/contest`;
  const response = await fetch(url);
  const result = await response.json();

  return result;
}

async function fetchProblems(contestId: string | null): Promise<string[]> {
  const url = contestId != null ? `${env.API_HOST}/api/list/problem?contest_id=${contestId}` : `${env.API_HOST}/api/list/problem`;
  const response = await fetch(url);
  const result = await response.json();

  return result;
}

export async function load({ url }): Promise<Data> {
  const filterCategory = url.searchParams.get("filter.category")?.split(",");
  let category: string | null = null;
  if (filterCategory != null) {
    category = filterCategory[0];
  }

  const filterContestId = url.searchParams.get("filter.contest_id")?.split(",");
  let contestId: string | null = null;
  if (filterContestId != null) {
    contestId = filterContestId[0];
  }

  const [result, categories, languages, contests, problems] = await Promise.all([fetchResult(url), fetchCategories(), fetchLanguages(), fetchContests(category), fetchProblems(contestId)]);

  return {
    result,
    categories: ["", ...categories],
    languages: ["", ...languages],
    contests: ["", ...contests],
    problems: ["", ...problems],
  };
}
