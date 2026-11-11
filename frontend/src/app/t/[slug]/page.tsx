import type { Metadata } from "next";

import { PublicDownloadView } from "@/features/transfers/PublicDownloadView";

export const metadata: Metadata = { title: "Download" };

type Props = { params: Promise<{ slug: string }> };

export default async function PublicDownloadPage({ params }: Props) {
  const { slug } = await params;
  return <PublicDownloadView slug={slug} />;
}
