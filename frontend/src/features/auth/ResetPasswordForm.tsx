"use client";

import Link from "next/link";
import { FormEvent, useState } from "react";

import { Alert } from "@/components/ui/Alert";
import { BrandLogo } from "@/components/ui/BrandLogo";
import { Button } from "@/components/ui/Button";
import { Icon } from "@/components/ui/Icon";
import { TextField } from "@/components/ui/TextField";
import { useRequestPasswordResetMutation } from "@/features/auth/authApi";
import { getApiErrorMessage } from "@/lib/apiError";

export function ResetPasswordForm() {
  const [requestReset, { isLoading, error, isSuccess }] =
    useRequestPasswordResetMutation();
  const [email, setEmail] = useState("");

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    try {
      await requestReset({ email: email.trim() }).unwrap();
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
            Reset Password
          </h1>
          <p className="mt-1 text-sm text-on-surface-variant">
            Enter your email and we&apos;ll send you a link to reset your
            password.
          </p>
        </div>
      </div>

      {error ? <Alert variant="error">{getApiErrorMessage(error)}</Alert> : null}
      {isSuccess ? (
        <Alert variant="success">
          <p className="font-semibold">Email sent</p>
          <p className="mt-1 text-on-surface">
            Please check your inbox for instructions to reset your password.
            (In development the token is logged by the API.)
          </p>
        </Alert>
      ) : null}

      <TextField
        label="Email address"
        name="email"
        type="email"
        autoComplete="email"
        required
        placeholder="you@example.com"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        leadingIcon={<Icon name="mail" />}
      />

      <Button type="submit" fullWidth loading={isLoading}>
        Send Reset Link
        <Icon name="arrow_forward" />
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
