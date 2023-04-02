import recoil from "recoil";
import { SearchResponse, FacetResults, Item } from "../types/response";

export const searchResponseState = recoil.atom<SearchResponse>({
  key: "searchResponseState",
  default: {
    stats: {
      total: 0,
      count: 0,
      pages: 0,
      index: 0,
      time: 0,
      facet: {
        category: {
          counts: [],
        },
        difficulty: {
          counts: [],
          start: "0",
          end: "4000",
          gap: "400",
          before: null,
          between: null,
          after: null,
        },
      },
    },
    items: [],
    message: null,
  },
});

/**
 * レスポンスから総ヒット数を取得するセレクタ
 */
export const searchResponseTotalSelector = recoil.selector<number>({
  key: "searchResponseTotalSelector",
  get: ({ get }) => {
    const response = get(searchResponseState);
    return response.stats.total;
  },
});

/**
 * レスポンスから1ページ当たりの表示件数を取得するセレクタ
 */
export const searchResponseCountSelector = recoil.selector<number>({
  key: "searchResponseCountSelector",
  get: ({ get }) => {
    const response = get(searchResponseState);
    return response.stats.count;
  },
});

/**
 * レスポンスから総ページ数を取得するセレクタ
 */
export const searchResponsePagesSelector = recoil.selector<number>({
  key: "searchResponsePagesSelector",
  get: ({ get }) => {
    const response = get(searchResponseState);
    return response.stats.pages;
  },
});

/**
 * レスポンスからページ番号を取得するセレクタ
 */
export const searchResponseIndexSelector = recoil.selector<number>({
  key: "searchResponseIndexSelector",
  get: ({ get }) => {
    const response = get(searchResponseState);
    return response.stats.index;
  },
});

/**
 * レスポンスから処理時間を取得するセレクタ
 */
export const searchResponseTimeSelector = recoil.selector<number>({
  key: "searchResponseTimeSelector",
  get: ({ get }) => {
    const response = get(searchResponseState);
    return response.stats.time;
  },
});

/**
 * レスポンスからファセット情報を取得するセレクタ
 */
export const searchResponseFacetSelector = recoil.selector<FacetResults>({
  key: "searchResponseFacetSelector",
  get: ({ get }) => {
    const response = get(searchResponseState);
    return response.stats.facet;
  },
});

/**
 * レスポンスからドキュメントを取得するセレクタ
 */
export const searchResponseItemsSelector = recoil.selector<Item[]>({
  key: "searchResponseItemsSelector",
  get: ({ get }) => {
    const response = get(searchResponseState);
    return response.items;
  },
});
