export type TransferStatus = "active" | "expired" | "deleted" | "blocked";

export type TransferFile = {
  id: string;
  original_name: string;
  mime_type: string;
  size_bytes: number;
  created_at: string;
};

export type Transfer = {
  id: string;
  short_code: string;
  title?: string | null;
  message?: string | null;
  status: TransferStatus | string;
  expires_at: string;
  download_count: number;
  has_password: boolean;
  share_url: string;
  created_at: string;
};

export type TransferListItem = Transfer & {
  file_count: number;
  total_size_bytes: number;
  recipient_count: number;
  first_file_name?: string | null;
};

export type TransferDetail = Transfer & {
  files: TransferFile[];
  file_count: number;
  total_size_bytes: number;
  recipient_count: number;
};

export type PublicTransfer = {
  short_code: string;
  title?: string | null;
  message?: string | null;
  status: TransferStatus | string;
  expires_at: string;
  has_password: boolean;
  sender_email: string;
  files: TransferFile[];
  file_count: number;
  total_size_bytes: number;
};

export type CreateTransferRequest = {
  title?: string;
  message?: string;
  password?: string;
  expires_in_hours: number;
  recipients?: string[];
};

export type TransferFilter = "active" | "expired" | "all";
