import axios from "axios";
import { CreateLinkRequest, Link } from "./types";

export const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

export async function getLinks(): Promise<Link[]> {
  const response = await api.get<Link[]>("/api/links");
  return response.data;
}

export async function createLink(request: CreateLinkRequest): Promise<Link> {
  const response = await api.post<Link>("/api/links", request);
  return response.data;
}
