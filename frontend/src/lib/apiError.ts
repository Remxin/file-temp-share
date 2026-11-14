import type { FetchBaseQueryError } from "@reduxjs/toolkit/query";
import type { SerializedError } from "@reduxjs/toolkit";

import type { ApiErrorBody } from "@/lib/types/auth";

export function getApiErrorMessage(
  error: FetchBaseQueryError | SerializedError | unknown,
  fallback = "Something went wrong",
): string {
  if (!error || typeof error !== "object") return fallback;
  if ("status" in error) {
    const fetchErr = error as FetchBaseQueryError;
    const data = fetchErr.data as ApiErrorBody | string | undefined;
    if (typeof data === "string" && data.trim()) return data;
    if (data && typeof data === "object" && "error" in data && data.error) {
      return data.error;
    }
    if (fetchErr.status === 401) return "Invalid email or password";
    if (fetchErr.status === 409) return "An account with this email already exists";
    if (fetchErr.status === 413) return "File is too large";
    if (fetchErr.status === "FETCH_ERROR") return "Cannot reach the API";
  }
  if ("message" in error && typeof (error as SerializedError).message === "string") {
    return (error as SerializedError).message || fallback;
  }
  return fallback;
}
