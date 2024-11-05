export const parseRange = (
  value: string | undefined
): { from?: number; to?: number } => {
  if (value == null || value === "") {
    return {};
  }

  if (!value.includes("-")) {
    return {};
  }

  let from = undefined;
  let to = undefined;

  const [fromStr, toStr] = value.split("-");
  if (fromStr !== "") {
    from = Number(fromStr);
  }
  if (toStr !== "") {
    to = Number(toStr);
  }

  return {
    from,
    to,
  };
};

export const parseBoolean = (
  value: string | undefined
): boolean | undefined => {
  switch (value) {
    case null:
      return undefined;
    case "true":
      return true;
    case "false":
      return false;
  }
};
