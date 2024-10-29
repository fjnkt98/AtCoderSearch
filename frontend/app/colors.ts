export const bgColorStyles = new Map<string, string>([
  ["ABC", "bg-blue-500"],
  ["ABC-Like", "bg-sky-600"],
  ["AGC", "bg-yellow-600"],
  ["AGC-Like", "bg-amber-500"],
  ["AHC", "bg-green-500"],
  ["ARC", "bg-red-500"],
  ["ARC-Like", "bg-orange-700"],
  ["JAG", "bg-slate-500"],
  ["JOI", "bg-slate-600"],
  ["Marathon", "bg-slate-600"],
  ["Other Contests", "bg-slate-600"],
  ["Other Sponsored", "bg-slate-600"],
  ["PAST", "bg-slate-600"],
]);

export const difficultyToBgColor = (difficulty: number | null): string => {
  if (difficulty == null) {
    return "bg-gray-900";
  }
  if (difficulty < 0) {
    return "bg-gray-900";
  }
  if (difficulty < 400) {
    return "bg-gray-500";
  }
  if (difficulty < 800) {
    return "bg-amber-800";
  }
  if (difficulty < 1200) {
    return "bg-green-500";
  }
  if (difficulty < 1600) {
    return "bg-cyan-400";
  }
  if (difficulty < 2000) {
    return "bg-blue-600";
  }
  if (difficulty < 2400) {
    return "bg-yellow-300";
  }
  if (difficulty < 2800) {
    return "bg-orange-400";
  }
  if (difficulty < 3200) {
    return "bg-red-600";
  }
  if (difficulty < 3600) {
    return "bg-zinc-400";
  }
  if (difficulty < 4000) {
    return "bg-amber-400";
  }
  return "";
};

export const difficultyToTextColor = (difficulty: number | null): string => {
  if (difficulty == null) {
    return "bg-gray-900";
  }
  if (difficulty < 0) {
    return "bg-gray-900";
  }
  if (difficulty < 400) {
    return "bg-gray-500";
  }
  if (difficulty < 800) {
    return "bg-amber-800";
  }
  if (difficulty < 1200) {
    return "bg-green-500";
  }
  if (difficulty < 1600) {
    return "bg-cyan-400";
  }
  if (difficulty < 2000) {
    return "bg-blue-600";
  }
  if (difficulty < 2400) {
    return "bg-yellow-300";
  }
  if (difficulty < 2800) {
    return "bg-orange-400";
  }
  if (difficulty < 3200) {
    return "bg-red-600";
  }
  if (difficulty < 3600) {
    return "bg-zinc-400";
  }
  if (difficulty < 4000) {
    return "bg-amber-400";
  }
  return "";
};

export const textColorStyles = new Map<string, string>([
  ["ABC", "text-blue-500"],
  ["ABC-Like", "text-sky-600"],
  ["AGC", "text-yellow-600"],
  ["AGC-Like", "text-amber-500"],
  ["AHC", "text-green-500"],
  ["ARC", "text-red-500"],
  ["ARC-Like", "text-orange-700"],
  ["JAG", "text-slate-500"],
  ["JOI", "text-slate-600"],
  ["Marathon", "text-slate-600"],
  ["Other Contests", "text-slate-600"],
  ["Other Sponsored", "text-slate-600"],
  ["PAST", "text-slate-600"],
  ["black", "text-gray-900"],
  ["gray", "text-gray-700"],
  ["brown", "text-amber-800"],
  ["green", "text-green-500"],
  ["cyan", "text-cyan-500"],
  ["blue", "text-blue-600"],
  ["yellow", "text-yellow-300"],
  ["orange", "text-orange-400"],
  ["red", "text-red-600"],
  ["silver", "text-zinc-600"],
  ["gold", "text-amber-500"],
]);
