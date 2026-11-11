export type TransferStatus = "active" | "expired";

export type TransferSummary = {
  id: string;
  title: string;
  status: TransferStatus;
  fileName: string;
  fileCount: number;
  totalSizeLabel: string;
  expiresAtLabel: string;
  downloadCount: number;
  shareUrl: string;
  shareMode: "public" | "recipients";
};

export type TransferFile = {
  id: string;
  name: string;
  sizeLabel: string;
};

export type TransferDetail = TransferSummary & {
  message?: string;
  uploadedAtLabel: string;
  files: TransferFile[];
};

export type PublicTransfer = {
  slug: string;
  title: string;
  senderName: string;
  message?: string;
  expiresAtLabel: string;
  status: TransferStatus;
  files: TransferFile[];
  totalSizeLabel: string;
};
