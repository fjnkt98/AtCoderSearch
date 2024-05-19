import type { SortValue } from "$lib/request";

export const selections = new Map<string, SortValue>([
  ["1", { label: "検索スコア順", values: ["-score"] }],
  ["2", { label: "レート低い順", values: ["rating"] }],
  ["3", { label: "レート高い順", values: ["-rating"] }],
  ["4", { label: "誕生年早い順", values: ["-birthYear"] }],
  ["5", { label: "誕生年遅い順", values: ["birthYear"] }],
]);
