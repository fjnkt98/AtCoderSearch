import { useEffect, useState } from "react";
import { SearchResponse, Item } from "./types/response";

export default function App() {
  const [items, setItems] = useState<Item[]>([]);

  useEffect(() => {
    (async () => {
      const api_endpoint = import.meta.env.VITE_API_ENDPOINT;

      const api_url = new URL(
        api_endpoint == null ? "http://localhost:8000/api/" : api_endpoint
      );
      const response = await fetch(api_url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          q: "高橋",
        }),
      });

      const content: SearchResponse = await response.json();
      console.log(content);
      setItems(content.items.docs);
    })();
  }, []);

  return (
    <>
      <p>Hello, world!</p>
      {items.map((item) => (
        <div key={item.problem_id}>
          <p>{item.problem_id}</p>
        </div>
      ))}
    </>
  );
}
