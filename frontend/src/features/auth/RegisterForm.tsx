"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

import { Alert } from "@/components/ui/Alert";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Checkbox } from "@/components/ui/Checkbox";
import { Icon } from "@/components/ui/Icon";
import { TextField } from "@/components/ui/TextField";
import { useRegisterMutation } from "@/features/auth/authApi";
import { getApiErrorMessage } from "@/lib/apiError";

export function RegisterForm() {
  const router = useRouter();
  const [register, { isLoading, error }] = useRegisterMutation();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [terms, setTerms] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
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
    if (!terms) {
      setLocalError("Please accept the Terms and Privacy Policy");
      return;
    }

    try {
      await register({ email: email.trim(), password }).unwrap();
      router.replace("/transfers");
    } catch {
      /* surfaced via error */
    }
  }

  return (
    <form className="flex flex-col gap-5" onSubmit={onSubmit}>
      <div className="flex flex-col items-center gap-3 text-center">
        <Badge dot>FlashShare</Badge>
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-on-surface sm:text-3xl">
            Create Account
          </h1>
          <p className="mt-1 text-sm text-on-surface-variant">
            Sign up to start sharing files
          </p>
        </div>
      </div>

      {(localError || error) && (
        <Alert variant="error">
          {localError ?? getApiErrorMessage(error)}
        </Alert>
      )}

      <TextField
        label="Email address"
        name="email"
        type="email"
        autoComplete="email"
        required
        placeholder="jane.doe@example.com"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        leadingIcon={<Icon name="mail" />}
      />

      <TextField
        label="Password"
        hint="Min. 8 characters"
        name="password"
        type={showPassword ? "text" : "password"}
        autoComplete="new-password"
        required
        minLength={8}
        placeholder="Enter password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        leadingIcon={<Icon name="lock" />}
        trailing={
          <button
            type="button"
            className="inline-flex size-10 items-center justify-center rounded-md text-on-surface-variant"
            aria-label={showPassword ? "Hide password" : "Show password"}
            onClick={() => setShowPassword((v) => !v)}
          >
            <Icon name={showPassword ? "visibility_off" : "visibility"} />
          </button>
        }
      />

      <TextField
        label="Confirm Password"
        name="confirm_password"
        type={showPassword ? "text" : "password"}
        autoComplete="new-password"
        required
        placeholder="Repeat password"
        value={confirm}
        onChange={(e) => setConfirm(e.target.value)}
        leadingIcon={<Icon name="lock_clock" />}
        error={
          confirm && password !== confirm ? "Passwords do not match." : undefined
        }
      />

      <Checkbox
        name="terms"
        checked={terms}
        onChange={(e) => setTerms(e.target.checked)}
        label={
          <>
            I agree to the{" "}
            <span className="font-semibold text-primary">Terms</span> and{" "}
            <span className="font-semibold text-primary">Privacy Policy</span>
          </>
        }
      />

      <Button type="submit" fullWidth loading={isLoading}>
        Create Account
        <Icon name="arrow_forward" />
      </Button>

      <p className="text-center text-sm text-on-surface-variant">
        Already have an account?{" "}
        <Link href="/login" className="font-semibold text-primary">
          Sign In
        </Link>
      </p>
    </form>
  );
}
