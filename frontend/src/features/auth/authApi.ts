import { baseApi } from "@/store/api/baseApi";
import {
  setCredentials,
  setUser,
} from "@/features/auth/authSlice";
import type {
  AuthResponse,
  ChangePasswordRequest,
  LoginRequest,
  MessageResponse,
  PasswordResetConfirmRequest,
  PasswordResetRequest,
  RegisterRequest,
  User,
} from "@/lib/types/auth";

export const authApi = baseApi.injectEndpoints({
  endpoints: (build) => ({
    register: build.mutation<AuthResponse, RegisterRequest>({
      query: (body) => ({
        url: "/api/v1/auth/register",
        method: "POST",
        body,
      }),
      async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          dispatch(setCredentials({ token: data.token, user: data.user }));
        } catch {
          /* form handles error */
        }
      },
    }),
    login: build.mutation<AuthResponse, LoginRequest>({
      query: (body) => ({
        url: "/api/v1/auth/login",
        method: "POST",
        body,
      }),
      async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          dispatch(setCredentials({ token: data.token, user: data.user }));
        } catch {
          /* form handles error */
        }
      },
    }),
    me: build.query<User, void>({
      query: () => "/api/v1/auth/me",
      providesTags: ["User"],
      async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          dispatch(setUser(data));
        } catch {
          /* hydrate / guard handles */
        }
      },
    }),
    changePassword: build.mutation<void, ChangePasswordRequest>({
      query: (body) => ({
        url: "/api/v1/auth/password",
        method: "POST",
        body,
      }),
    }),
    requestPasswordReset: build.mutation<MessageResponse, PasswordResetRequest>({
      query: (body) => ({
        url: "/api/v1/auth/password-reset",
        method: "POST",
        body,
      }),
    }),
    confirmPasswordReset: build.mutation<void, PasswordResetConfirmRequest>({
      query: (body) => ({
        url: "/api/v1/auth/password-reset/confirm",
        method: "POST",
        body,
      }),
    }),
  }),
});

export const {
  useRegisterMutation,
  useLoginMutation,
  useMeQuery,
  useLazyMeQuery,
  useChangePasswordMutation,
  useRequestPasswordResetMutation,
  useConfirmPasswordResetMutation,
} = authApi;
