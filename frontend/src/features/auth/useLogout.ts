"use client";

import { useRouter } from "next/navigation";

import { clearCredentials } from "@/features/auth/authSlice";
import { baseApi } from "@/store/api/baseApi";
import { useAppDispatch } from "@/store/hooks";

export function useLogout() {
  const dispatch = useAppDispatch();
  const router = useRouter();

  return () => {
    dispatch(clearCredentials());
    dispatch(baseApi.util.resetApiState());
    router.replace("/login");
  };
}
