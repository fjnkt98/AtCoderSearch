export const intoURLSearchParams = (object: object): URLSearchParams => {
  const result = new URLSearchParams();
  Object.entries(object).forEach(([key, value]) => {
    if (value == null) {
      return;
    }

    if (Array.isArray(value)) {
      value.forEach((v) => result.append(key, String(v)));
    } else {
      result.set(key, value);
    }
  });
  return result;
};

export const numberFromQueryString = (s: string | null): number | null => {
  if (s == null) {
    return null;
  }
  if (s === "") {
    return null;
  }

  return Number(s);
};

export const nullableBooleanFromQueryString = (s: string | null): boolean | null => {
  if (s == null) {
    return null;
  }

  if (s === "true") {
    return true;
  } else if (s === "false") {
    return false;
  } else {
    return null;
  }
};

export const booleanFromQueryString = (s: string | null): boolean | null => {
  if (s == null) {
    return null;
  }
  if (s === "true") {
    return true;
  } else {
    return null;
  }
};

export function parseRange(label: string): NumericRange {
  const [begin, end] = label.split("~");
  const result: NumericRange = { begin: null, end: null };
  if (begin.trim() !== "") {
    result.begin = Number(begin.trim());
  }
  if (end.trim() !== "") {
    result.end = Number(end.trim());
  }
  return result;
}

export function rangeFromQueryStrings(params: URLSearchParams, fromKey: string, toKey: string): string | null {
  const from = params.get(fromKey);
  const to = params.get(toKey);

  if (from == null && to == null) {
    return null;
  }
  const fromValue = from ?? "";
  const toValue = to ?? "";

  return `${fromValue} ~ ${toValue}`;
}

export interface SortValue {
  label: string;
  values: string[];
}

export interface NumericRange {
  begin: number | null;
  end: number | null;
}

export interface SearchProblemParameter {
  limit: number | null;
  page: number | null;
  q: string | null;
  sort: string[] | null;
  facet: string[] | null;
  category: string[] | null;
  difficultyFrom: number | null;
  difficultyTo: number | null;
  color: string[] | null;
  userId: string | null;
  difficulty: number | null;
  excludeSolved: boolean | null;
  experimental: boolean | null;
  prioritizeRecent: boolean | null;
}

export interface RecommendProblemParameter {
  limit: number | null;
  page: number | null;
  problemId: string | null;
}

export interface SearchUserParameter {
  limit: number | null;
  page: number | null;
  q: string | null;
  sort: string[] | null;
  facet: string[] | null;
  userId: string[] | null;
  ratingFrom: number | null;
  ratingTo: number | null;
  birthYearFrom: number | null;
  birthYearTo: number | null;
  joinCountFrom: number | null;
  joinCountTo: number | null;
  country: string[] | null;
  color: string[] | null;
}

export interface SearchSubmissionParameter {
  limit: number | null;
  page: number | null;
  sort: string[] | null;
  epochSecondFrom: number | null;
  epochSecondTo: number | null;
  problemId: string[] | null;
  contestId: string[] | null;
  category: string[] | null;
  userId: string[] | null;
  language: string[] | null;
  languageGroup: string[] | null;
  pointFrom: number | null;
  pointTo: number | null;
  lengthFrom: number | null;
  lengthTo: number | null;
  result: string[] | null;
  executionTimeFrom: number | null;
  executionTimeTo: number | null;
}
