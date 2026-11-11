import type { Metadata } from "next";

import { AuthShell } from "@/components/layout/AuthShell";
import { GuestOnly } from "@/features/auth/GuestOnly";
import { LoginForm } from "@/features/auth/LoginForm";

export const metadata: Metadata = { title: "Sign In" };

export default function LoginPage() {
  return (
    <GuestOnly>
      <AuthShell>
        <LoginForm />
      </AuthShell>
    </GuestOnly>
  );
}
