"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { FormEvent, Suspense, useState } from "react";

import { Alert } from "@/components/ui/Alert";
import { BrandLogo } from "@/components/ui/BrandLogo";
import { Button } from "@/components/ui/Button";
import { Icon } from "@/components/ui/Icon";
import { Spinner } from "@/components/ui/Spinner";
import { TextField } from "@/components/ui/TextField";
import { useConfirmPasswordResetMutation } from "@/features/auth/authApi";
import { getApiErrorMessage } from "@/lib/apiError";

function ConfirmResetFormInner() {
  const router = useRouter();
  const params = useSearchParams();
  const tokenFromQuery = params.get("token") ?? "";
  const [confirmReset, { isLoading, error, isSuccess }] =
    useConfirmPasswordResetMutation();
  const [token, setToken] = useState(tokenFromQuery);
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [localError, setLocalError] = useState<string | null>(null);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setLocalError(null);
    if (password.length < 8) {
      setLocalError("Password must be at least 8 characters");
      return;
    }
    if (password !== confirm) {
      setLocalError("Passwords do not match.");
      return;
    }
    try {
      await confirmReset({ token: token.trim(), new_password: password }).unwrap();
      setTimeout(() => router.replace("/login"), 1200);
    } catch {
      /* surfaced */
    }
  }

  return (
    <form className="flex flex-col gap-5" onSubmit={onSubmit}>
      <div className="flex flex-col items-center gap-3 text-center">
        <BrandLogo href={null} size="sm" />
        <div>
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">
            Choose a new password
          </h1>
          <p className="mt-1 text-sm text-on-surface-variant">
            Paste the reset token from your email (or logs) and set a new
            password.
          </p>
        </div>
      </div>

      {(localError || error) && (
        <Alert variant="error">
          {localError ?? getApiErrorMessage(error)}
        </Alert>
      )}
      {isSuccess ? (
        <Alert variant="success">Password updated. Redirecting to Sign In…</Alert>
      ) : null}

      <TextField
        label="Reset token"
        name="token"
        required
        placeholder="Paste token"
        value={token}
        onChange={(e) => setToken(e.target.value)}
        leadingIcon={<Icon name="vpn_key" />}
      />

      <TextField
        label="New password"
        hint="Min. 8 characters"
        name="new_password"
        type="password"
        autoComplete="new-password"
        required
        minLength={8}
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        leadingIcon={<Icon name="lock" />}
      />

      <TextField
        label="Confirm password"
        name="confirm_password"
        type="password"
        autoComplete="new-password"
        required
        value={confirm}
        onChange={(e) => setConfirm(e.target.value)}
        leadingIcon={<Icon name="lock_clock" />}
      />

      <Button type="submit" fullWidth loading={isLoading}>
        Update password
      </Button>

      <Link
        href="/login"
        className="inline-flex min-h-[var(--touch-min)] items-center justify-center gap-1 text-sm font-semibold text-primary"
      >
        <Icon name="arrow_back" />
        Back to Sign In
      </Link>
    </form>
  );
}

export function ConfirmResetForm() {
  return (
    <Suspense
      fallback={
        <div className="flex justify-center py-10">
          <Spinner className="text-3xl" />
        </div>
      }
    >
      <ConfirmResetFormInner />
    </Suspense>
  );
}
