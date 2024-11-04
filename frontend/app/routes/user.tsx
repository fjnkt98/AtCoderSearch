import { useLoaderData, json, Form, useSearchParams } from "@remix-run/react";
import { useState } from "react";
import { client } from "~/client";
import { z } from "zod";
import { zx } from "zodix";
import type { LoaderFunctionArgs } from "@remix-run/node";
import Pagination from "~/components/Pagination";
import { parseRange } from "~/utils";
import { difficultyToTextColor } from "~/colors";

export const loader = async ({ request }: LoaderFunctionArgs) => {
  const { q, page, sort, rating, birthYear, joinCount, country } =
    zx.parseQuery(request, {
      q: z.string().optional(),
      page: zx.IntAsString.optional(),
      sort: z.string().optional(),
      rating: z.string().optional(),
      birthYear: z.string().optional(),
      joinCount: z.string().optional(),
      country: z
        .union([z.string().array(), z.string()])
        .transform((v) => (Array.isArray(v) ? v : [v]))
        .optional()
        .default([]),
    });

  const parseSort = (
    sort: string | undefined
  ): (
    | "rating:asc"
    | "rating:desc"
    | "birthYear:asc"
    | "birthYear:desc"
    | "joinCount:asc"
    | "joinCount:desc"
    | "rank:asc"
    | "rank:desc"
    | "accepted:asc"
    | "accepted:desc"
    | "submissionCount:asc"
    | "submissionCount:desc"
  )[] => {
    switch (sort) {
      case "-rating":
        return ["rating:desc", "rank:asc"];
      case "rating":
        return ["rating:asc", "rank:desc"];
      case "-birthYear":
        return ["birthYear:desc", "rank:asc"];
      case "birthYear":
        return ["birthYear:asc", "rank:asc"];
      case "-joinCount":
        return ["joinCount:desc", "rank:asc"];
      case "joinCount":
        return ["joinCount:asc", "rank:asc"];
      case "-accepted":
        return ["accepted:desc", "rank:asc"];
      case "accepted":
        return ["accepted:asc", "rank:asc"];
      case "-submissionCount":
        return ["submissionCount:desc", "rank:asc"];
      case "submissionCount":
        return ["submissionCount:asc", "rank:asc"];
      default:
        return ["rating:desc", "rank:asc"];
    }
  };

  const { data, error } = await client.POST("/api/user", {
    body: {
      q: q,
      limit: 50,
      page: page ?? 1,
      sort: parseSort(sort),
      facet: ["country", "rating", "birthYear", "joinCount"],
      rating: parseRange(rating),
      birthYear: parseRange(birthYear),
      joinCount: parseRange(joinCount),
      country: country,
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
  const [ratingLastSelected, setRatingLastSelected] = useState("");
  const [birthYearLastSelected, setBirthYearLastSelected] = useState("");
  const [joinCountLastSelected, setJoinCountLastSelected] = useState("");

  const ratings = data.facet.rating ?? [];
  const birthYears = data.facet.birthYear ?? [];
  const joinCounts = data.facet.joinCount ?? [];
  const countries = data.facet.country ?? [];

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
              <option value="-rating">レート高い順</option>
              <option value="rating">レート低い順</option>
              <option value="-birthYear">誕生年早い順</option>
              <option value="birthYear">誕生年遅い順</option>
              <option value="-joinCount">参加回数多い順</option>
              <option value="joinCount">参加回数少ない順</option>
              <option value="-accepted">AC数多い順</option>
              <option value="accepted">AC数少ない順</option>
              <option value="-submissionCount">総提出回数多い順</option>
              <option value="submissionCount">総提出回数少ない順</option>
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
                <p className="text-lg">レーティング</p>
                <div className="flex flex-wrap gap-0.5">
                  {ratings.map((c) => (
                    <label key={c.label} className="block p-0.5 select-none">
                      <input
                        id={`rating${c.label
                          .replaceAll(" ", "")
                          .replaceAll("~", "-")}`}
                        type="radio"
                        name="rating"
                        className="mr-1"
                        value={c.label.replaceAll(" ", "").replaceAll("~", "-")}
                        onClick={(e) => {
                          if (e.currentTarget.id === ratingLastSelected) {
                            e.currentTarget.checked = !e.currentTarget.checked;
                            setRatingLastSelected("");
                          } else {
                            setRatingLastSelected(e.currentTarget.id);
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
                <p className="text-lg">誕生年</p>
                <div className="flex flex-wrap gap-0.5">
                  {birthYears.map((c) => (
                    <label key={c.label} className="block p-0.5 select-none">
                      <input
                        id={`birthYear${c.label
                          .replaceAll(" ", "")
                          .replaceAll("~", "-")}`}
                        type="radio"
                        name="birthYear"
                        className="mr-1"
                        value={c.label.replaceAll(" ", "").replaceAll("~", "-")}
                        onClick={(e) => {
                          if (e.currentTarget.id === birthYearLastSelected) {
                            e.currentTarget.checked = !e.currentTarget.checked;
                            setBirthYearLastSelected("");
                          } else {
                            setBirthYearLastSelected(e.currentTarget.id);
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
                <p className="text-lg">参加回数</p>
                <div className="flex flex-wrap gap-0.5">
                  {joinCounts.map((c) => (
                    <label key={c.label} className="block p-0.5 select-none">
                      <input
                        id={`joinCount${c.label
                          .replaceAll(" ", "")
                          .replaceAll("~", "-")}`}
                        type="radio"
                        name="joinCount"
                        className="mr-1"
                        value={c.label.replaceAll(" ", "").replaceAll("~", "-")}
                        onClick={(e) => {
                          if (e.currentTarget.id === joinCountLastSelected) {
                            e.currentTarget.checked = !e.currentTarget.checked;
                            setJoinCountLastSelected("");
                          } else {
                            setJoinCountLastSelected(e.currentTarget.id);
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
                <p className="text-lg">国</p>
                <div className="flex flex-wrap gap-0.5">
                  {countries.map((c) => (
                    <label key={c.label} className="block p-0.5 select-none">
                      <input
                        type="checkbox"
                        name="country"
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
            </div>

            <input type="hidden" name="page" value="1" />

            <div className="flex flex-row items-center justify-end">
              <button
                type="submit"
                className="px-3 py-2 rounded-lg bg-blue-600 text-gray-100"
              >
                絞り込む
              </button>
            </div>
          </Form>
        </div>

        <div className="overflow-x-auto">
          <table className="table-auto text-sm">
            <thead>
              <th scope="col" className="border px-2 py-1">
                順位
              </th>
              <th scope="col" className="border px-2 py-1">
                ユーザID
              </th>
              <th scope="col" className="border px-2 py-1">
                所属
              </th>
              <th scope="col" className="border px-2 py-1">
                誕生年
              </th>
              <th scope="col" className="border px-2 py-1">
                レート
              </th>
              <th scope="col" className="border px-2 py-1">
                参加回数
              </th>
              <th scope="col" className="border px-2 py-1">
                総提出回数
              </th>
              <th scope="col" className="border px-2 py-1">
                AC数
              </th>
            </thead>
            <tbody>
              {data.items.map((item) => (
                <tr key={item.userId} className="">
                  <td className="text-center border px-2 py-1">{item.rank}</td>
                  <td className="text-center border px-2 py-1">
                    <a
                      href={item.userUrl}
                      rel="noreferrer"
                      target="_blank"
                      className={difficultyToTextColor(item.rating)}
                    >
                      {item.userId}
                    </a>
                  </td>
                  <td className="text-center border truncate px-2 py-1">
                    {item.affiliation}
                  </td>
                  <td className="text-center border px-2 py-1">
                    {item.birthYear}
                  </td>
                  <td className="text-center border px-2 py-1">
                    {item.rating}
                  </td>
                  <td className="text-center border px-2 py-1">
                    {item.joinCount}
                  </td>
                  <td className="text-center border px-2 py-1">
                    {item.submissionCount}
                  </td>
                  <td className="text-center border px-2 py-1">
                    {item.accepted}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <div className="py-2">
        <Pagination
          to="/user"
          params={searchParams}
          current={data.index}
          end={data.pages}
        />
      </div>
    </>
  );
}
