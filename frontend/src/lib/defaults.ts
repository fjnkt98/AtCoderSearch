export const defaultProblemSearchParams = new URLSearchParams([
  ["limit", "20"],
  ["page", "1"],
  ["facet.term", "category"],
  ["facet.difficulty.from", "0"],
  ["facet.difficulty.to", "4000"],
  ["facet.difficulty.gap", "400"],
]);

export const defaultUserSearchParams = new URLSearchParams([
  ["limit", "50"],
  ["page", "1"],
  ["facet.term", "country"],
  ["facet.rating.from", "0"],
  ["facet.rating.to", "4000"],
  ["facet.rating.gap", "400"],
  ["facet.birth_year.from", "1970"],
  ["facet.birth_year.to", "2020"],
  ["facet.birth_year.gap", "10"],
  ["facet.join_count.from", "0"],
  ["facet.join_count.to", "100"],
  ["facet.join_count.gap", "20"],
]);

export const defaultSubmissionSearchParams = new URLSearchParams([
  ["limit", "50"],
  ["page", "1"],
]);
