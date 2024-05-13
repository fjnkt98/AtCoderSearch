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
