import { useLoaderData, json, Form, useSearchParams } from "@remix-run/react";
import { useState } from "react";
import { client } from "~/client";
import { z } from "zod";
import { zx } from "zodix";
import type { LoaderFunctionArgs } from "@remix-run/node";
import { categoryToTextColor, difficultyToTextColor } from "~/colors";
import Pagination from "~/components/Pagination";
import { parseBoolean, parseRange } from "~/utils";

export const loader = async ({ request }: LoaderFunctionArgs) => {
  const { q, page, sort, category, difficulty, userId, experimental } =
    zx.parseQuery(request, {
      q: z.string().optional(),
      page: zx.IntAsString.optional(),
      sort: z.string().optional(),
      category: z
        .union([z.string().array(), z.string()])
        .transform((v) => (Array.isArray(v) ? v : [v]))
        .optional()
        .default([]),
      difficulty: z.string().optional(),
      userId: z.string().optional(),
      experimental: z.string().optional(),
    });

  const parseSort = (
    sort: string | undefined
  ): (
    | "startAt:asc"
    | "startAt:desc"
    | "difficulty:asc"
    | "difficulty:desc"
    | "problemId:asc"
    | "problemId:desc"
  )[] => {
    switch (sort) {
      case "-startAt":
        return ["startAt:desc"];
      case "startAt":
        return ["startAt:asc"];
      case "-difficulty":
        return ["difficulty:desc"];
      case "difficulty":
        return ["difficulty:asc"];
      default:
        return [];
    }
  };

  const { data, error } = await client.POST("/api/problem", {
    body: {
      q: q,
      limit: 50,
      page: page ?? 1,
      sort: parseSort(sort),
      facet: ["category", "difficulty"],
      category: category,
      difficulty: parseRange(difficulty),
      experimental: parseBoolean(experimental),
      userId,
    },
  });

  if (error != null) {
    throw new Response("request failed", {
      status: 500,
      statusText: "INTERNAL SERVER ERROR",
    });
  }

  return json(data);
};

export default function ProblemPage() {
  const data = useLoaderData<typeof loader>();
  const [searchParams] = useSearchParams();
  const [menuOpen, setMenuOpen] = useState(false);
  const [difficultyLastSelected, setDifficultyLastSelected] = useState("");
  const [experimentalLastSelected, setExperimentalLastSelected] = useState("");

  const categories = data.facet.category ?? [];
  const difficulties = data.facet.difficulty ?? [];

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
        <span>
          約{data.total}件 / {data.time}ms
        </span>
      </div>

      <div className="lg:flex lg:flex-row">
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
              <option value="-startAt">新しい順</option>
              <option value="startAt">古い順</option>
              <option value="-difficulty">難易度高い順</option>
              <option value="difficulty">難易度低い順</option>
            </select>

            <div className="flex flex-col gap-3">
              <label className="block">
                <span className="block text-lg">キーワード</span>
                <input
                  type="text"
                  name="q"
                  className="border rounded-md border-gray-400 shadow-sm"
                />
              </label>

              <div>
                <p className="text-lg">カテゴリ</p>
                <div className="flex flex-wrap gap-0.5">
                  {categories.map((c) => (
                    <label key={c.label} className="block p-0.5 select-none">
                      <input
                        type="checkbox"
                        name="category"
                        className="mr-1"
                        value={c.label}
                      />
                      <span className="">
                        {c.label}({c.count})
                      </span>
                    </label>
                  ))}
                </div>
              </div>

              <div>
                <p className="text-lg">難易度</p>
                <div className="flex flex-wrap gap-0.5">
                  {difficulties.map((c) => (
                    <label key={c.label} className="block p-0.5 select-none">
                      <input
                        type="radio"
                        id={c.label}
                        name="difficulty"
                        className="mr-1"
                        value={c.label.replaceAll(" ", "").replaceAll("~", "-")}
                        onClick={(e) => {
                          if (e.currentTarget.id === difficultyLastSelected) {
                            e.currentTarget.checked = !e.currentTarget.checked;
                            setDifficultyLastSelected("");
                          } else {
                            setDifficultyLastSelected(e.currentTarget.id);
                          }
                        }}
                      />
                      <span className="">
                        {c.label}({c.count})
                      </span>
                    </label>
                  ))}
                </div>
              </div>

              <div>
                <p className="text-lg">その他</p>

                <div className="flex flex-col gap-2">
                  <label className="block">
                    <span className="block">解いたことのある問題を除外</span>
                    <input
                      type="text"
                      placeholder="ユーザID"
                      name="userId"
                      className="border rounded-md border-gray-400 shadow-sm px-2"
                    />
                  </label>

                  <div className="flex flex-row gap-2">
                    <label>
                      <input
                        id="experimentalFalse"
                        type="radio"
                        name="experimental"
                        className="mr-1"
                        value="false"
                        onClick={(e) => {
                          if (e.currentTarget.id === experimentalLastSelected) {
                            e.currentTarget.checked = !e.currentTarget.checked;
                            setExperimentalLastSelected("");
                          } else {
                            setExperimentalLastSelected(e.currentTarget.id);
                          }
                        }}
                      />
                      <span>試験管問題を除く</span>
                    </label>

                    <label>
                      <input
                        id="experimentalTrue"
                        type="radio"
                        name="experimental"
                        className="mr-1"
                        value="true"
                        onClick={(e) => {
                          if (e.currentTarget.id === experimentalLastSelected) {
                            e.currentTarget.checked = !e.currentTarget.checked;
                            setExperimentalLastSelected("");
                          } else {
                            setExperimentalLastSelected(e.currentTarget.id);
                          }
                        }}
                      />
                      <span>試験管問題のみ</span>
                    </label>
                  </div>
                </div>
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

        <div className="flex flex-col lg:flex-grow items-center px-2 text-sm gap-1 py-1">
          {data.items.map((item) => (
            <div
              key={item.problemId}
              className="border border-gray-500 dark:border-gray-400 rounded-md min-w-60 max-w-xl w-full py-2 flex flex-row items-center justify-between gap-2 px-2"
            >
              <a
                href={item.contestUrl}
                target="_blank"
                rel="noreferrer"
                className={`text-center ${categoryToTextColor(item.category)}`}
              >
                {item.contestId}
              </a>
              <a
                href={item.problemUrl}
                target="_blank"
                rel="noreferrer"
                className={`text-balance text-center ${difficultyToTextColor(
                  item.difficulty
                )}`}
              >
                {item.problemTitle}
              </a>
              <span className="min-w-8">{item.difficulty}</span>
            </div>
          ))}
        </div>
      </div>
      <Pagination
        to="/problem"
        params={searchParams}
        current={data.index}
        end={data.pages}
      />
    </>
  );
}
