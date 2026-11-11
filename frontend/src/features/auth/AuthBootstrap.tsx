"use client";

import { useEffect } from "react";

import {
  finishHydrationAnonymous,
  startHydration,
} from "@/features/auth/authSlice";
import { useLazyMeQuery } from "@/features/auth/authApi";
import { selectAuthStatus } from "@/features/auth/selectors";
import { Spinner } from "@/components/ui/Spinner";
import { getStoredToken } from "@/lib/storage/token";
import { useAppDispatch, useAppSelector } from "@/store/hooks";

export function AuthBootstrap({ children }: { children: React.ReactNode }) {
  const dispatch = useAppDispatch();
  const status = useAppSelector(selectAuthStatus);
  const [fetchMe] = useLazyMeQuery();

  useEffect(() => {
    let cancelled = false;

    async function hydrate() {
      dispatch(startHydration());
      const token = getStoredToken();
      if (!token) {
        if (!cancelled) dispatch(finishHydrationAnonymous());
        return;
      }
      try {
        await fetchMe().unwrap();
      } catch {
        if (!cancelled) dispatch(finishHydrationAnonymous());
      }
    }

    if (status === "idle") {
      void hydrate();
    }

    return () => {
      cancelled = true;
    };
  }, [dispatch, fetchMe, status]);

  if (status === "idle" || status === "hydrating") {
    return (
      <div className="flex min-h-dvh items-center justify-center bg-surface">
        <Spinner className="text-3xl" />
      </div>
    );
  }

  return children;
}
