"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { Spinner } from "@/components/ui/Spinner";
import {
  selectCurrentUser,
  selectIsAuthReady,
  selectIsAuthenticated,
} from "@/features/auth/selectors";
import { useAppSelector } from "@/store/hooks";

export function RequireAdmin({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const ready = useAppSelector(selectIsAuthReady);
  const authenticated = useAppSelector(selectIsAuthenticated);
  const user = useAppSelector(selectCurrentUser);
  const isAdmin = user?.role === "ADMIN";

  useEffect(() => {
    if (!ready) return;
    if (!authenticated) {
      router.replace("/login");
      return;
    }
    if (!isAdmin) {
      router.replace("/transfers");
    }
  }, [ready, authenticated, isAdmin, router]);

  if (!ready || !authenticated || !isAdmin) {
    return (
      <div className="flex min-h-dvh items-center justify-center">
        <Spinner className="text-3xl" />
      </div>
    );
  }

  return children;
}
