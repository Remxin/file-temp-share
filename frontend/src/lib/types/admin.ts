import type { TransferListItem } from "@/lib/types/transfer";

export type AdminTransferListItem = TransferListItem & {
  user_id: string;
  owner_email: string;
};

export type AdminSetUserActiveRequest = {
  is_active: boolean;
};

export type AdminSetTransferStatusRequest = {
  status: "blocked" | "active";
};
