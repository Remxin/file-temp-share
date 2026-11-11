import type { Metadata } from "next";

import { AuthShell } from "@/components/layout/AuthShell";
import { GuestOnly } from "@/features/auth/GuestOnly";
import { RegisterForm } from "@/features/auth/RegisterForm";

export const metadata: Metadata = { title: "Create Account" };

export default function RegisterPage() {
  return (
    <GuestOnly>
      <AuthShell>
        <RegisterForm />
      </AuthShell>
    </GuestOnly>
  );
}
