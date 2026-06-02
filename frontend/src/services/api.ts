import type { ApiEnvelope } from "../types/api";

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? "";
const TOKEN_KEY = "fail2ban-dashboard-token";

export function getToken() {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
}

// Global filters reference accessed by fetch wrapper
export const currentFilters = {
  domainId: 0,
  startTime: "",
  endTime: "",
};

export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
  const token = getToken();
  const headers = new Headers(init.headers);
  headers.set("Content-Type", "application/json");
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  // Auto-inject filters into GET requests
  const method = (init.method || "GET").toUpperCase();
  if (method === "GET") {
    // Avoid applying filters on loading the domains list itself to prevent recursion
    if (!path.includes("/api/domains")) {
      const url = new URL(path, "http://localhost"); // dummy host to parse relative path
      if (currentFilters.domainId > 0) {
        url.searchParams.set("domain_id", String(currentFilters.domainId));
      }
      if (currentFilters.startTime) {
        url.searchParams.set("start_time", currentFilters.startTime);
      }
      if (currentFilters.endTime) {
        url.searchParams.set("end_time", currentFilters.endTime);
      }
      path = url.pathname + url.search;
    }
  }

  const res = await fetch(`${API_BASE}${path}`, { ...init, headers });
  const payload = (await res.json()) as ApiEnvelope<T>;
  if (!res.ok || !payload.success) {
    throw new Error(payload.error || "Request failed");
  }
  return payload.data as T;
}
