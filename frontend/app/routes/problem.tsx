import { useLoaderData, json } from "@remix-run/react";
import createClient from "openapi-fetch";
import type { paths } from "~/types";

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
    <div>
      <p>Search Problem</p>
      {data.items.map((item) => (
        <div key={item.problemId}>
          <p>{item.problemId}</p>
        </div>
      ))}
    </div>
  );
}
