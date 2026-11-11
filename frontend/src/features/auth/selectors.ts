import type { RootState } from "@/store";

export const selectAuthStatus = (state: RootState) => state.auth.status;
export const selectAuthToken = (state: RootState) => state.auth.token;
export const selectCurrentUser = (state: RootState) => state.auth.user;
export const selectIsAuthenticated = (state: RootState) =>
  state.auth.status === "authenticated";
export const selectIsAuthReady = (state: RootState) =>
  state.auth.status === "authenticated" || state.auth.status === "anonymous";
