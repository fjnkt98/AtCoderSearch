import { useLoaderData, json, Form } from "@remix-run/react";
import { useState } from "react";
import { client } from "~/client";

export const loader = async () => {
  const { data, error } = await client.POST("/api/problem", {
    body: {
      facet: ["category", "difficulty"],
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
  const [menuOpen, setMenuOpen] = useState(false);
  const [lastSelected, setLastSelected] = useState("");

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
          } lg:block absolute lg:static w-5/6 lg:w-1/4 bg-white dark:bg-gray-950 px-4 py-2 shadow-md shadow-gray-500 rounded-xl`}
        >
          <Form method="GET">
            <div className="flex flex-row items-center justify-between text-lg">
              <span className="text-lg font-bold">絞り込み</span>
              <button type="button" className="lg:hidden" onClick={toggle}>
                ✗
              </button>
            </div>

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
                          if (e.currentTarget.id === lastSelected) {
                            e.currentTarget.checked = !e.currentTarget.checked;
                            setLastSelected("");
                          } else {
                            setLastSelected(e.currentTarget.id);
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

                  <label className="block">
                    <input
                      type="checkbox"
                      name="isExperimental"
                      className="mr-1"
                    />
                    <span>試験管問題を除外</span>
                  </label>
                </div>
              </div>
            </div>

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

        <div className="flex flex-col lg:flex-grow items-center px-2 text-sm gap-1 py-1">
          {data.items.map((item) => (
            <div
              key={item.problemId}
              className="border rounded-md min-w-60 max-w-xl w-full py-2 flex flex-row items-center justify-between gap-2 px-2"
            >
              <a href={item.contestUrl} className="text-center">
                {item.contestId}
              </a>
              <a href={item.problemUrl} className="text-balance text-center">
                {item.problemTitle}
              </a>
              <span className="min-w-8">{item.difficulty}</span>
            </div>
          ))}
        </div>
      </div>
    </>
  );
}
