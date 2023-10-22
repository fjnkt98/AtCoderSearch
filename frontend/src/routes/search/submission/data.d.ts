import type { SubmissionResult } from "$lib/search";

export type Data = {
  result: SubmissionResult;
  categories: string[];
  languages: string[];
  languageGroups: string[];
  contests: string[];
  problems: string[];
};
