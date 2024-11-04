import {
  useRouteError,
  isRouteErrorResponse,
  useLoaderData,
  json,
  Form,
  useSearchParams,
} from "@remix-run/react";
import { useState } from "react";
import { client } from "~/client";
import { z } from "zod";
import { zx } from "zodix";
import dayjs from "dayjs";
import timezone from "dayjs/plugin/timezone";
import utc from "dayjs/plugin/utc";
import type { LoaderFunctionArgs } from "@remix-run/node";
import Pagination from "~/components/Pagination";
import {
  categoryToTextColor,
  difficultyToTextColor,
  resultToTextColor,
} from "~/colors";

dayjs.extend(utc);
dayjs.extend(timezone);

const jst = "Asia/Tokyo";

export const loader = async ({ request }: LoaderFunctionArgs) => {
  const {
    page,
    sort,
    epochSecondFrom,
    epochSecondTo,
    problemId,
    contestId,
    category,
    userId,
    languageGroup,
    pointFrom,
    pointTo,
    lengthFrom,
    lengthTo,
    result,
    executionTimeFrom,
    executionTimeTo,
  } = zx.parseQuery(request, {
    page: zx.IntAsString.optional(),
    sort: z.string().optional(),
    epochSecondFrom: z.string().optional(),
    epochSecondTo: z.string().optional(),
    problemId: z
      .union([z.string().array(), z.string()])
      .transform((v) => (Array.isArray(v) ? v : [v]))
      .optional()
      .default([]),
    contestId: z
      .union([z.string().array(), z.string()])
      .transform((v) => (Array.isArray(v) ? v : [v]))
      .optional()
      .default([]),
    category: z
      .union([z.string().array(), z.string()])
      .transform((v) => (Array.isArray(v) ? v : [v]))
      .optional()
      .default([]),
    userId: z
      .union([z.string().array(), z.string()])
      .transform((v) => (Array.isArray(v) ? v : [v]))
      .optional()
      .default([]),
    languageGroup: z
      .union([z.string().array(), z.string()])
      .transform((v) => (Array.isArray(v) ? v : [v]))
      .optional()
      .default([]),
    pointFrom: z.string().optional(),
    pointTo: z.string().optional(),
    lengthFrom: z.string().optional(),
    lengthTo: z.string().optional(),
    result: z
      .union([z.string().array(), z.string()])
      .transform((v) => (Array.isArray(v) ? v : [v]))
      .optional()
      .default([]),
    executionTimeFrom: z.string().optional(),
    executionTimeTo: z.string().optional(),
  });

  const parseSort = (
    sort: string | undefined
  ): (
    | "executionTime:asc"
    | "executionTime:desc"
    | "epochSecond:asc"
    | "epochSecond:desc"
    | "point:asc"
    | "point:desc"
    | "length:asc"
    | "length:desc"
  )[] => {
    switch (sort) {
      case "-executionTime":
        return ["executionTime:desc"];
      case "executionTime":
        return ["executionTime:asc"];
      case "-epochSecond":
        return ["epochSecond:desc"];
      case "epochSecond":
        return ["epochSecond:asc"];
      case "-point":
        return ["point:desc"];
      case "point":
        return ["point:asc"];
      case "-length":
        return ["length:desc"];
      case "length":
        return ["length:asc"];
      default:
        return ["epochSecond:desc"];
    }
  };

  const parseNumString = (value: string | undefined): number | undefined => {
    if (value == null) {
      return undefined;
    }

    if (value === "") {
      return undefined;
    }

    const num = Number(value);
    if (Number.isNaN(num)) {
      return undefined;
    }
    return num;
  };

  const parseLocalDateTime = (
    value: string | undefined
  ): number | undefined => {
    if (value == null) {
      return undefined;
    }

    if (value === "") {
      return undefined;
    }

    const dt = dayjs(value, "YYYY-MM-DDTHH:mm").tz(jst);
    const timestamp = dt.unix();

    if (Number.isNaN(timestamp)) {
      return undefined;
    }
    return timestamp;
  };

  const [categoryRes, contestRes, languageRes, problemRes, submissionRes] =
    await Promise.all([
      client.GET("/api/category"),
      client.GET("/api/contest", {
        params: { query: { category: category.filter((e) => e !== "") } },
      }),
      client.GET("/api/language", {
        params: { query: { group: languageGroup.filter((e) => e !== "") } },
      }),
      client.GET("/api/problem", {
        params: {
          query: {
            category: category.filter((e) => e !== ""),
            contestId: contestId.filter((e) => e !== ""),
          },
        },
      }),
      client.POST("/api/submission", {
        body: {
          limit: 50,
          page: page ?? 1,
          sort: parseSort(sort),
          problemId: problemId.filter((e) => e !== ""),
          contestId: contestId.filter((e) => e !== ""),
          category: category.filter((e) => e !== ""),
          userId: userId.filter((e) => e !== ""),
          languageGroup: languageGroup.filter((e) => e !== ""),
          result: result.filter((e) => e !== ""),
          epochSecond: {
            from: parseLocalDateTime(epochSecondFrom),
            to: parseLocalDateTime(epochSecondTo),
          },
          point: {
            from: parseNumString(pointFrom),
            to: parseNumString(pointTo),
          },
          length: {
            from: parseNumString(lengthFrom),
            to: parseNumString(lengthTo),
          },
          executionTime: {
            from: parseNumString(executionTimeFrom),
            to: parseNumString(executionTimeTo),
          },
        },
      }),
    ]);

  const errors = [
    categoryRes.error,
    contestRes.error,
    languageRes.error,
    problemRes.error,
    submissionRes.error,
  ].filter((e) => e != null);

  if (errors.length > 0) {
    console.error(errors);
    throw new Response(errors.map((e) => e.message).join(": "), {
      status: 500,
      statusText: "INTERNAL SERVER ERROR",
    });
  }

  return json({
    category: categoryRes.data,
    contest: contestRes.data,
    language: languageRes.data,
    problem: problemRes.data,
    submission: submissionRes.data,
  });
};

export function ErrorBoundary() {
  const error = useRouteError();

  if (isRouteErrorResponse(error)) {
    return (
      <div>
        <h1>
          {error.status} {error.statusText}
        </h1>
      </div>
    );
  }

  if (error instanceof Error) {
    <div>
      <h1>Error</h1>
      <p>{error.message}</p>
      <p>The stack trace is:</p>
      <pre>{error.stack}</pre>
    </div>;
  }

  return (
    <div>
      <h1>Unknown error</h1>
    </div>
  );
}

export default function SubmissionPage() {
  const data = useLoaderData<typeof loader>();
  const [searchParams] = useSearchParams();
  const [menuOpen, setMenuOpen] = useState(false);

  const toggle = () => {
    setMenuOpen(!menuOpen);
  };

  return (
    <>
      <div className="flex flex-row justify-between lg:justify-end items-center py-2 px-2">
        <button
          type="button"
          className={`lg:hidden rounded-xl px-3 py-1 border shadow-sm shadow-gray-500 ${
            menuOpen ? "bg-gray-200 dark:bg-gray-950 dark:text-gray-400" : ""
          }`}
          onClick={toggle}
        >
          絞り込み
        </button>
        <span>{data.submission?.time}ms</span>
      </div>

      <div className="flex flex-row">
        <div
          className={`${
            menuOpen ? "block" : "hidden"
          } lg:block absolute lg:static w-5/6 lg:w-1/4 bg-white dark:bg-gray-950 px-4 py-2 border border-gray-700 dark:border-gray-300 rounded-xl`}
        >
          <Form method="GET">
            <div className="flex flex-row items-center justify-between text-lg">
              <span className="text-lg font-bold">絞り込み</span>
              <button type="button" className="lg:hidden" onClick={toggle}>
                ✗
              </button>
            </div>

            <select name="sort">
              <option value="-epochSecond">提出日時新しい順</option>
              <option value="epochSecond">提出日時古い順</option>
              <option value="-executionTime">実行時間長い順</option>
              <option value="executionTime">実行時間短い順</option>
              <option value="-length">コード長長い順</option>
              <option value="length">コード長短い順</option>
              <option value="-point">得点高い順</option>
              <option value="point">得点低い順</option>
            </select>

            <div className="flex flex-row flex-wrap gap-4">
              <div>
                <p className="text-lg">カテゴリ</p>
                <select name="category">
                  <option value="">カテゴリ</option>
                  {data.category?.categories?.map((c) => (
                    <option key={c} value={c}>
                      {c}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <p className="text-lg">言語</p>
                <select name="languageGroup">
                  <option value="">言語</option>
                  {data.language?.languages?.map((l) => (
                    <option key={l.group} value={l.group}>
                      {l.group}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <p className="text-lg">コンテスト</p>
                <select name="contestId">
                  <option value="">コンテスト</option>
                  {data.contest?.contests?.map((c) => (
                    <option key={c} value={c}>
                      {c}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <p className="text-lg">問題</p>
                <select name="problemId">
                  <option value="">問題</option>
                  {data.problem?.problems?.map((p) => (
                    <option key={p} value={p}>
                      {p}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <p className="text-lg">結果</p>
                <select name="result">
                  <option value="" className="text-gray-500">
                    結果
                  </option>
                  <option value="AC">AC</option>
                  <option value="WA">WA</option>
                  <option value="TLE">TLE</option>
                  <option value="RE">RE</option>
                  <option value="CE">CE</option>
                  <option value="MLE">MLE</option>
                </select>
              </div>

              <div>
                <p className="text-lg">ユーザID</p>
                <input type="text" name="userId" placeholder="ユーザID" />
              </div>

              <div>
                <p className="text-lg">得点</p>
                <input
                  type="number"
                  name="pointFrom"
                  min={0}
                  placeholder="得点(下限)"
                  className="m-1"
                />
                <input
                  type="number"
                  name="pointTo"
                  min={0}
                  placeholder="得点(上限)"
                  className="m-1"
                />
              </div>

              <div>
                <p className="text-lg">コード長</p>
                <input
                  type="number"
                  name="lengthFrom"
                  min={0}
                  placeholder="コード長(下限)"
                  className="m-1"
                />
                <input
                  type="number"
                  name="lengthTo"
                  min={0}
                  placeholder="コード長(上限)"
                  className="m-1"
                />
              </div>

              <div>
                <p className="text-lg">実行時間</p>
                <input
                  type="number"
                  name="executionTimeFrom"
                  min={0}
                  placeholder="実行時間(下限)"
                  className="m-1"
                />
                <input
                  type="number"
                  name="executionTimeTo"
                  min={0}
                  placeholder="実行時間(上限)"
                  className="m-1"
                />
              </div>

              <div>
                <p className="text-lg">提出日時(JST)</p>
                <label className="block m-1">
                  <span className="mr-2">開始</span>
                  <input type="datetime-local" name="epochSecondFrom" />
                </label>
                <label className="block m-1">
                  <span className="mr-2">終了</span>
                  <input type="datetime-local" name="epochSecondTo" />
                </label>
              </div>
            </div>

            <input type="hidden" name="page" value="1" />

            <div className="flex flex-row items-center justify-end gap-4">
              <button
                type="submit"
                className="px-3 py-2 rounded-lg bg-blue-600 text-gray-100"
              >
                絞り込む
              </button>

              <input
                type="reset"
                value="リセット"
                className="block text-sm px-3 py-2 border border-gray-700 dark:border-gray-400 cursor-pointer select-none rounded-lg"
              />
            </div>
          </Form>
        </div>

        <div className="overflow-x-auto overflow-y-hidden">
          <table className="table-auto">
            <thead className="text-sm lg:text-md">
              <th
                scope="col"
                className="border border-gray-700 dark:border-gray-300 px-2 py-1"
              >
                ID
              </th>
              <th
                scope="col"
                className="border border-gray-700 dark:border-gray-300 px-2 py-1"
              >
                提出日時
              </th>
              <th
                scope="col"
                className="border border-gray-700 dark:border-gray-300 px-2 py-1"
              >
                コンテストID
              </th>
              <th
                scope="col"
                className="border border-gray-700 dark:border-gray-300 px-2 py-1"
              >
                タイトル
              </th>
              <th
                scope="col"
                className="border border-gray-700 dark:border-gray-300 px-2 py-1"
              >
                ユーザID
              </th>
              <th
                scope="col"
                className="border border-gray-700 dark:border-gray-300 px-2 py-1"
              >
                コード長
              </th>
              <th
                scope="col"
                className="border border-gray-700 dark:border-gray-300 px-2 py-1"
              >
                実行時間
              </th>
              <th
                scope="col"
                className="border border-gray-700 dark:border-gray-300 px-2 py-1"
              >
                結果
              </th>
              <th
                scope="col"
                className="border border-gray-700 dark:border-gray-300 px-2 py-1"
              >
                言語
              </th>
            </thead>
            <tbody className="text-xs lg:text-sm">
              {data.submission?.items.map((item) => (
                <tr key={item.submissionId} className="">
                  <td className="text-center border border-gray-700 dark:border-gray-300 px-2 py-1 text-bold text-blue-500 underline">
                    <a
                      href={item.submissionUrl}
                      rel="noreferrer"
                      target="_blank"
                    >
                      {item.submissionId}
                    </a>
                  </td>
                  <td className="text-center border border-gray-700 dark:border-gray-300 px-2 py-1 text-gray-700 dark:text-gray-300">
                    {dayjs(item.submittedAt * 1000)
                      .tz(jst)
                      .format("YYYY-MM-DDTHH:mm:ssZZ")}
                  </td>
                  <td className="text-center border border-gray-700 dark:border-gray-300 px-2 py-1">
                    <a
                      href={item.contestUrl}
                      rel="noreferrer"
                      target="_blank"
                      className={`${categoryToTextColor(item.category)}`}
                    >
                      {item.contestId}
                    </a>
                  </td>
                  <td className="text-center border border-gray-700 dark:border-gray-300 truncate px-2 py-1">
                    <a
                      href={item.problemUrl}
                      rel="noreferrer"
                      target="_blank"
                      className={`${difficultyToTextColor(
                        item.difficulty
                      )} text-wrap`}
                    >
                      {item.problemTitle}
                    </a>
                  </td>
                  <td className="text-center border border-gray-700 dark:border-gray-300 px-2 py-1">
                    {item.userId}
                  </td>
                  <td className="text-center border border-gray-700 dark:border-gray-300 px-2 py-1 text-wrap">
                    {item.length}
                  </td>
                  <td className="text-center border border-gray-700 dark:border-gray-300 px-2 py-1 text-wrap">
                    {item.point}
                  </td>
                  <td
                    className={`text-center border border-gray-700 dark:border-gray-300 px-2 py-1 ${resultToTextColor(
                      item.result
                    )}`}
                  >
                    {item.result}
                  </td>
                  <td className="text-center border border-gray-700 dark:border-gray-300 px-2 py-1 text-wrap">
                    {item.language}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <div className="py-2">
        <Pagination
          to="/submission"
          params={searchParams}
          current={data.submission?.index ?? 1}
        />
      </div>
    </>
  );
}
