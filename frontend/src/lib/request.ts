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
