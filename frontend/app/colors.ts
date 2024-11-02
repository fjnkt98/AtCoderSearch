export const categoryToTextColor = (category: string): string => {
  switch (category) {
    case "ABC":
    case "ABC-Like":
      return "text-blue-600 dark:text-blue-400";
    case "AGC":
    case "AGC-Like":
      return "text-amber-500 dark:text-yellow-500";
    case "AHC":
    case "Marathon":
      return "text-green-500 dark:text-green-500";
    case "ARC":
    case "ARC-Like":
      return "text-red-500 dark:text-red-500";
    case "JAG":
    case "JOI":
    case "Other Contests":
    case "Other Sponsored":
      return "text-slate-700 dark:text-slate-200";
    case "PAST":
      return "text-pink-500 dark:text-pink-400";
    default:
      return "";
  }
};

export const difficultyToTextColor = (
  difficulty: number | undefined
): string => {
  if (difficulty == null) {
    return "text-gray-950 dark:text-gray-100";
  }
  if (difficulty < 0) {
    return "text-gray-950 dark:text-gray-100";
  }
  if (difficulty < 400) {
    return "text-gray-600 dark:text-gray-400";
  }
  if (difficulty < 800) {
    return "text-amber-800 dark:text-amber-700";
  }
  if (difficulty < 1200) {
    return "text-green-600 dark:text-green-500";
  }
  if (difficulty < 1600) {
    return "text-cyan-500 dark:text-cyan-400";
  }
  if (difficulty < 2000) {
    return "text-blue-600 dark:text-blue-500";
  }
  if (difficulty < 2400) {
    return "text-yellow-500 dark:text-yellow-400";
  }
  if (difficulty < 2800) {
    return "text-orange-500 dark:text-orange-400";
  }
  if (difficulty < 3200) {
    return "text-red-600 dark:text-red-500";
  }
  if (difficulty < 3600) {
    return "text-zinc-400 dark:text-zinc-300";
  }
  return "text-amber-500 dark:text-amber-300";
};
