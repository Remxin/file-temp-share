import type { Metadata } from "next";

import { TransferDetailView } from "@/features/transfers/TransferDetailView";

export const metadata: Metadata = { title: "Transfer details" };

type Props = { params: Promise<{ id: string }> };

export default async function TransferDetailPage({ params }: Props) {
  const { id } = await params;
  return <TransferDetailView id={id} />;
}
