"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { Spinner } from "@/components/ui/Spinner";
import {
  selectIsAuthReady,
  selectIsAuthenticated,
} from "@/features/auth/selectors";
import { useAppSelector } from "@/store/hooks";

export function RequireAuth({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const ready = useAppSelector(selectIsAuthReady);
  const authenticated = useAppSelector(selectIsAuthenticated);

  useEffect(() => {
    if (ready && !authenticated) {
      router.replace("/login");
    }
  }, [ready, authenticated, router]);

  if (!ready || !authenticated) {
    return (
      <div className="flex min-h-dvh items-center justify-center">
        <Spinner className="text-3xl" />
      </div>
    );
  }

  return children;
}
