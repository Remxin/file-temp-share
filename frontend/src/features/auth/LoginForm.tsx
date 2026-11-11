"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

import { Alert } from "@/components/ui/Alert";
import { Button } from "@/components/ui/Button";
import { Checkbox } from "@/components/ui/Checkbox";
import { Icon } from "@/components/ui/Icon";
import { TextField } from "@/components/ui/TextField";
import { useLoginMutation } from "@/features/auth/authApi";
import { getApiErrorMessage } from "@/lib/apiError";

export function LoginForm() {
  const router = useRouter();
  const [login, { isLoading, error }] = useLoginMutation();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [remember, setRemember] = useState(true);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    try {
      await login({ email: email.trim(), password }).unwrap();
      // remember is UX-only for now (token always persisted for session restore)
      void remember;
      router.replace("/transfers");
    } catch {
      /* surfaced via error */
    }
  }

  return (
    <form className="flex flex-col gap-5" onSubmit={onSubmit}>
      <div className="flex flex-col items-center gap-3 text-center">
        <span className="inline-flex size-12 items-center justify-center rounded-xl bg-primary-container text-on-primary shadow-md">
          <Icon name="cyclone" className="text-[1.6rem]" />
        </span>
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-on-surface sm:text-3xl">
            Sign In
          </h1>
          <p className="mt-1 text-sm text-on-surface-variant">
            Enter your email and password to access your transfers
          </p>
        </div>
      </div>

      {error ? <Alert variant="error">{getApiErrorMessage(error)}</Alert> : null}

      <TextField
        label="Email address"
        name="email"
        type="email"
        autoComplete="email"
        required
        placeholder="name@company.com"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        leadingIcon={<Icon name="mail" />}
      />

      <TextField
        label="Password"
        name="password"
        type={showPassword ? "text" : "password"}
        autoComplete="current-password"
        required
        placeholder="••••••••"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        leadingIcon={<Icon name="key" />}
        trailing={
          <button
            type="button"
            className="inline-flex size-10 items-center justify-center rounded-md text-on-surface-variant hover:text-on-surface"
            aria-label={showPassword ? "Hide password" : "Show password"}
            onClick={() => setShowPassword((v) => !v)}
          >
            <Icon name={showPassword ? "visibility_off" : "visibility"} />
          </button>
        }
      />

      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <Checkbox
          name="remember_me"
          checked={remember}
          onChange={(e) => setRemember(e.target.checked)}
          label="Remember me"
        />
        <Link
          href="/reset-password"
          className="min-h-[var(--touch-min)] text-sm font-semibold text-primary inline-flex items-center"
        >
          Forgot password?
        </Link>
      </div>

      <Button type="submit" fullWidth loading={isLoading}>
        Sign In
      </Button>

      <p className="text-center text-sm text-on-surface-variant">
        Don&apos;t have an account?{" "}
        <Link href="/register" className="font-semibold text-primary">
          Create an account
        </Link>
      </p>
    </form>
  );
}
