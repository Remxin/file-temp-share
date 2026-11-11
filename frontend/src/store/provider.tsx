"use client";

import { useRef } from "react";
import { Provider } from "react-redux";

import { AuthBootstrap } from "@/features/auth/AuthBootstrap";
import { makeStore, type AppStore } from "@/store";

export function StoreProvider({ children }: { children: React.ReactNode }) {
  const storeRef = useRef<AppStore | null>(null);
  if (!storeRef.current) {
    storeRef.current = makeStore();
  }

  return (
    <Provider store={storeRef.current}>
      <AuthBootstrap>{children}</AuthBootstrap>
    </Provider>
  );
}
