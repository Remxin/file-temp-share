import type { Metadata } from "next";

import { Card } from "@/components/ui/Card";

export const metadata: Metadata = { title: "Settings" };

export default function SettingsPage() {
  return (
    <Card className="max-w-xl space-y-2">
      <h1 className="text-2xl font-bold">Settings</h1>
      <p className="text-sm text-on-surface-variant">
        Account settings (email / password) will live here. Stub for now.
      </p>
    </Card>
  );
}
