"use client";

import { FormEvent, useState } from "react";

import { Alert } from "@/components/ui/Alert";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Icon } from "@/components/ui/Icon";
import { TextField } from "@/components/ui/TextField";
import { useChangePasswordMutation } from "@/features/auth/authApi";
import { getApiErrorMessage } from "@/lib/apiError";

export function ChangePasswordForm() {
  const [changePassword, { isLoading, error, isSuccess, reset }] =
    useChangePasswordMutation();
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [localError, setLocalError] = useState<string | null>(null);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setLocalError(null);
    reset();

    if (newPassword.length < 8) {
      setLocalError("New password must be at least 8 characters");
      return;
    }
    if (newPassword !== confirmPassword) {
      setLocalError("New passwords do not match.");
      return;
    }
    if (currentPassword === newPassword) {
      setLocalError("New password must differ from the current password.");
      return;
    }

    try {
      await changePassword({
        current_password: currentPassword,
        new_password: newPassword,
      }).unwrap();
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
    } catch {
      /* surfaced via error */
    }
  }

  return (
    <Card>
      <form className="space-y-5" onSubmit={onSubmit}>
        <div>
          <h2 className="text-base font-bold tracking-tight sm:text-lg">
            Change password
          </h2>
          <p className="mt-1 text-sm text-on-surface-variant">
            Enter your current password, then choose a new one.
          </p>
        </div>

        {(localError || error) && (
          <Alert variant="error">
            {localError ??
              getApiErrorMessage(error, "Could not change password")}
          </Alert>
        )}
        {isSuccess ? (
          <Alert variant="success">Password updated successfully.</Alert>
        ) : null}

        <div className="grid gap-4 sm:grid-cols-1 lg:grid-cols-3">
          <TextField
            label="Current password"
            name="current_password"
            type="password"
            autoComplete="current-password"
            required
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
            leadingIcon={<Icon name="lock" />}
          />
          <TextField
            label="New password"
            hint="Min. 8 characters"
            name="new_password"
            type="password"
            autoComplete="new-password"
            required
            minLength={8}
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            leadingIcon={<Icon name="key" />}
          />
          <TextField
            label="Confirm new password"
            name="confirm_password"
            type="password"
            autoComplete="new-password"
            required
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            leadingIcon={<Icon name="lock_clock" />}
          />
        </div>

        <Button type="submit" loading={isLoading} className="w-full sm:w-auto">
          Update password
        </Button>
      </form>
    </Card>
  );
}
