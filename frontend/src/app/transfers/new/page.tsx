import type { Metadata } from "next";

import { NewTransferWizard } from "@/features/transfers/NewTransferWizard";

export const metadata: Metadata = { title: "New Transfer" };

export default function NewTransferPage() {
  return <NewTransferWizard />;
}
