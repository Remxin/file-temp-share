import type { Metadata } from "next";

import { AccountOverview } from "@/features/settings/AccountOverview";

export const metadata: Metadata = { title: "Settings" };

export default function SettingsPage() {
  return <AccountOverview />;
}
