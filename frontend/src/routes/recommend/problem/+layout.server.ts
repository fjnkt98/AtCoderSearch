export const csr = false;

export async function load({ url, cookies }) {
  let userId: string | null = url.searchParams.get("user_id");
  if (userId == null) {
    userId = cookies.get("user_id") ?? null;
  } else {
    cookies.set("user_id", userId, { path: "/recommend/problem" });
  }

  return {
    userId: userId,
  };
}
