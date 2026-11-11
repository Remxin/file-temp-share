import type { Metadata } from "next";

import { AuthShell } from "@/components/layout/AuthShell";
import { GuestOnly } from "@/features/auth/GuestOnly";
import { ConfirmResetForm } from "@/features/auth/ConfirmResetForm";

export const metadata: Metadata = { title: "Set new password" };

export default function ConfirmResetPasswordPage() {
  return (
    <GuestOnly>
      <AuthShell>
        <ConfirmResetForm />
      </AuthShell>
    </GuestOnly>
  );
}
