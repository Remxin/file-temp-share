import type { Metadata } from "next";

import { AuthShell } from "@/components/layout/AuthShell";
import { GuestOnly } from "@/features/auth/GuestOnly";
import { ResetPasswordForm } from "@/features/auth/ResetPasswordForm";

export const metadata: Metadata = { title: "Reset Password" };

export default function ResetPasswordPage() {
  return (
    <GuestOnly>
      <AuthShell>
        <ResetPasswordForm />
      </AuthShell>
    </GuestOnly>
  );
}
