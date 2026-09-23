import axios from "axios";
import { RegisterShortcutRequest, Shortcut } from "./types";

export const BACKEND_ORIGIN =
  import.meta.env.VITE_BACKEND_ORIGIN ?? "http://localhost:8080";

const client = axios.create({
  baseURL: BACKEND_ORIGIN,
  headers: {
    "Content-Type": "application/json",
  },
});

export async function fetchShortcuts(): Promise<Shortcut[]> {
  const response = await client.get<Shortcut[]>("/api/shortcuts");
  return response.data;
}

export async function registerShortcut(
  request: RegisterShortcutRequest,
): Promise<Shortcut> {
  const response = await client.post<Shortcut>("/api/shortcuts", request);
  return response.data;
}
