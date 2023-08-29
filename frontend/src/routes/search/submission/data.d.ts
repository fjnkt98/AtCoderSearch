import type { SubmissionResult } from "$lib/search";

export type Data = {
  result: SubmissionResult;
  categories: string[];
  languages: string[];
  contests: string[];
  problems: string[];
};
