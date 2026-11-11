import type {
  PublicTransfer,
  TransferDetail,
  TransferSummary,
} from "@/lib/types/transfer";

export const mockTransfers: TransferSummary[] = [
  {
    id: "t1",
    title: "Client_Brand_Guidelines_2025.pdf",
    status: "active",
    fileName: "Client_Brand_Guidelines_2025.pdf",
    fileCount: 1,
    totalSizeLabel: "42 MB",
    expiresAtLabel: "Expires in 4 days (May 28, 2025)",
    downloadCount: 14,
    shareUrl: "https://flashshare.local/t/brand-2025",
    shareMode: "public",
  },
  {
    id: "t2",
    title: "Q2_Launch_Reel.mp4",
    status: "active",
    fileName: "Q2_Launch_Reel.mp4",
    fileCount: 1,
    totalSizeLabel: "860 MB",
    expiresAtLabel: "Expires in 1 day (May 25, 2025)",
    downloadCount: 3,
    shareUrl: "https://flashshare.local/t/q2-reel",
    shareMode: "public",
  },
  {
    id: "t3",
    title: "Archive_Spring.zip",
    status: "expired",
    fileName: "Archive_Spring.zip",
    fileCount: 1,
    totalSizeLabel: "120 MB",
    expiresAtLabel: "Expired May 20, 2025",
    downloadCount: 9,
    shareUrl: "https://flashshare.local/t/spring-archive",
    shareMode: "recipients",
  },
];

export const mockTransferDetails: Record<string, TransferDetail> = {
  t1: {
    ...mockTransfers[0],
    message:
      "Hi team, here are the finalized high-resolution brand assets for Q2 rollout.",
    uploadedAtLabel: "Uploaded May 20, 2025",
    files: [
      {
        id: "f1",
        name: "01_brand_guidelines_final.pdf",
        sizeLabel: "42 MB",
      },
    ],
  },
  t2: {
    ...mockTransfers[1],
    uploadedAtLabel: "Uploaded May 22, 2025",
    files: [
      { id: "f2", name: "Q2_Launch_Reel.mp4", sizeLabel: "860 MB" },
    ],
  },
  t3: {
    ...mockTransfers[2],
    uploadedAtLabel: "Uploaded May 10, 2025",
    files: [{ id: "f3", name: "Archive_Spring.zip", sizeLabel: "120 MB" }],
  },
};

export const mockPublicTransfers: Record<string, PublicTransfer> = {
  "brand-2025": {
    slug: "brand-2025",
    title: "Client_Brand_Guidelines_2025.pdf",
    senderName: "Alex Mercer",
    message: "Final brand guidelines for Q2.",
    expiresAtLabel: "Available until May 28, 2025",
    status: "active",
    totalSizeLabel: "42 MB",
    files: [
      {
        id: "f1",
        name: "Client_Brand_Guidelines_2025.pdf",
        sizeLabel: "42 MB",
      },
    ],
  },
  expired: {
    slug: "expired",
    title: "Client_Brand_Assets_2025.zip",
    senderName: "Alex Mercer",
    expiresAtLabel: "Expired on May 20, 2025",
    status: "expired",
    totalSizeLabel: "210 MB",
    files: [
      {
        id: "fx",
        name: "Client_Brand_Assets_2025.zip",
        sizeLabel: "210 MB",
      },
    ],
  },
};
