import type { SortValue } from "$lib/request";

export const selections = new Map<string, SortValue>([
  ["1", { label: "検索スコア順", values: ["-score", "-startAt"] }],
  ["2", { label: "新しい順", values: ["-startAt"] }],
  ["3", { label: "古い順", values: ["startAt"] }],
  ["4", { label: "難易度高い順", values: ["-difficulty", "-startAt"] }],
  ["5", { label: "難易度低い順", values: ["difficulty", "-startAt"] }],
]);
