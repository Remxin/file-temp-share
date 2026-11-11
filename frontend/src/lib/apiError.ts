import type { FetchBaseQueryError } from "@reduxjs/toolkit/query";
import type { SerializedError } from "@reduxjs/toolkit";

import type { ApiErrorBody } from "@/lib/types/auth";

export function getApiErrorMessage(
  error: FetchBaseQueryError | SerializedError | undefined,
  fallback = "Something went wrong",
): string {
  if (!error) return fallback;
  if ("status" in error) {
    const data = error.data as ApiErrorBody | string | undefined;
    if (typeof data === "string" && data.trim()) return data;
    if (data && typeof data === "object" && "error" in data && data.error) {
      return data.error;
    }
    if (error.status === 401) return "Invalid email or password";
    if (error.status === 409) return "An account with this email already exists";
    if (error.status === "FETCH_ERROR") return "Cannot reach the API";
  }
  if ("message" in error && error.message) return error.message;
  return fallback;
}
