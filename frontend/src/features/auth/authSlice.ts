import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

import type { User } from "@/lib/types/auth";
import {
  clearStoredToken,
  getStoredToken,
  setStoredToken,
} from "@/lib/storage/token";

export type AuthStatus = "idle" | "hydrating" | "authenticated" | "anonymous";

type AuthState = {
  token: string | null;
  user: User | null;
  status: AuthStatus;
};

const initialState: AuthState = {
  token: null,
  user: null,
  status: "idle",
};

type CredentialsPayload = {
  token: string;
  user: User;
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCredentials(state, action: PayloadAction<CredentialsPayload>) {
      state.token = action.payload.token;
      state.user = action.payload.user;
      state.status = "authenticated";
      setStoredToken(action.payload.token);
    },
    setUser(state, action: PayloadAction<User>) {
      state.user = action.payload;
      if (state.token) state.status = "authenticated";
    },
    clearCredentials(state) {
      state.token = null;
      state.user = null;
      state.status = "anonymous";
      clearStoredToken();
    },
    startHydration(state) {
      state.status = "hydrating";
      state.token = getStoredToken();
    },
    finishHydrationAnonymous(state) {
      state.token = null;
      state.user = null;
      state.status = "anonymous";
      clearStoredToken();
    },
  },
});

export const {
  setCredentials,
  setUser,
  clearCredentials,
  startHydration,
  finishHydrationAnonymous,
} = authSlice.actions;

export default authSlice.reducer;
