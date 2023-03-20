import recoil from "recoil";
import { createSearchParams } from "react-router-dom";

export const searchParamsState = recoil.atom<URLSearchParams>({
  key: "searchParamsState",
  default: createSearchParams({
    limit: "20",
    page: "1",
    sort: "-score",
  }),
});

export const searchParamsStateSelector = recoil.selector<URLSearchParams>({
  key: "searchParamsStateSelector",
  get: ({ get }) => {
    return get(searchParamsState);
  },
});
