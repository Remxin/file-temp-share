import type { Metadata } from "next";

import { TransfersDashboard } from "@/features/transfers/TransfersDashboard";

export const metadata: Metadata = { title: "Dashboard" };

export default function TransfersPage() {
  return <TransfersDashboard />;
}
