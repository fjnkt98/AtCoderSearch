import type { SortValue } from "$lib/request";

export const selections = new Map<string, SortValue>([
  ["1", { label: "提出日時早い順", values: ["-submittedAt"] }],
  ["2", { label: "提出日時遅い順", values: ["submittedAt"] }],
  ["3", { label: "実行時間短い順", values: ["executionTime"] }],
  ["4", { label: "実行時間長い順", values: ["-executionTime"] }],
  ["5", { label: "得点小さい順", values: ["point"] }],
  ["6", { label: "得点多い順", values: ["-point"] }],
  ["7", { label: "コード長短い順", values: ["length"] }],
  ["8", { label: "コード長長い順", values: ["-length"] }],
]);
