import { useLoaderData, json } from "@remix-run/react";
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

  return (
    <div className="">
      <div className="flex flex-row justify-between items-center py-2 px-2">
        <button
          type="button"
          className="rounded-xl px-3 py-1 border shadow-sm shadow-gray-500"
        >
          絞り込み
        </button>
        <span>約1,000件 / 30ms</span>
      </div>
      <div className="flex flex-col items-center px-2 text-sm gap-1 py-1">
        {data.items.map((item) => (
          <div
            key={item.problemId}
            className="border rounded-md min-w-60 w-full py-2 flex flex-row items-center justify-between gap-2 px-2"
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
