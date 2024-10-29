import { useLoaderData, json, Form } from "@remix-run/react";
import { useState } from "react";
import createClient from "openapi-fetch";
import type { paths } from "~/types";
// import {
//   Listbox,
//   ListboxButton,
//   ListboxOption,
//   ListboxOptions,
// } from "@headlessui/react";

const client = createClient<paths>({ baseUrl: "http://localhost:8000" });

export const loader = async () => {
  const { data, error } = await client.POST("/api/problem", {
    body: {},
  });

  if (error != null) {
    throw new Response("request failed");
  }

  return json(data);
};

export default function ProblemPage() {
  const data = useLoaderData<typeof loader>();
  const [menuOpen, setMenuOpen] = useState(false);

  const toggle = () => {
    setMenuOpen(!menuOpen);
  };

  return (
    <div className="">
      <div
        className={`${
          menuOpen ? "block" : "hidden"
        } absolute w-5/6 border border-blue-500 bg-white dark:bg-gray-950 px-4 py-2`}
      >
        <Form method="GET">
          <div className="flex flex-row items-center justify-between text-lg">
            <span>絞り込み</span>
            <button type="button" onClick={toggle}>
              ✗
            </button>
          </div>
          <p>キーワード</p>
          <input type="text" name="q" />

          <button type="submit">絞り込む</button>
          <button type="button">クリア</button>
        </Form>
      </div>
      <div className="flex flex-row justify-between lg:justify-end items-center py-2 px-2">
        <button
          type="button"
          className="lg:hidden rounded-xl px-3 py-1 border shadow-sm shadow-gray-500"
          onClick={toggle}
        >
          絞り込み
        </button>
        <span>
          約{data.total}件 / {data.time}ms
        </span>
      </div>
      <div className="flex flex-col items-center px-2 text-sm gap-1 py-1">
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
  );
}
