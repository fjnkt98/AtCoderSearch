import createClient from "openapi-fetch";
import type { paths } from "~/types";

export const client = createClient<paths>({ baseUrl: process.env.BACKEND_URL });
