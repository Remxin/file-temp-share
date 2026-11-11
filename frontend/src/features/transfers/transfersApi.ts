import { baseApi } from "@/store/api/baseApi";
import type {
  PublicTransfer,
  TransferDetail,
  TransferSummary,
} from "@/lib/types/transfer";
import {
  mockPublicTransfers,
  mockTransferDetails,
  mockTransfers,
} from "@/features/transfers/mockData";

export type TransferFilter = "active" | "expired" | "all";

export const transfersApi = baseApi.injectEndpoints({
  endpoints: (build) => ({
    listTransfers: build.query<TransferSummary[], TransferFilter>({
      async queryFn(filter) {
        await delay(200);
        const items =
          filter === "all"
            ? mockTransfers
            : mockTransfers.filter((t) => t.status === filter);
        return { data: items };
      },
      providesTags: ["Transfers"],
    }),
    getTransfer: build.query<TransferDetail, string>({
      async queryFn(id) {
        await delay(150);
        const item = mockTransferDetails[id];
        if (!item) return { error: { status: 404, data: "not found" } };
        return { data: item };
      },
      providesTags: (_r, _e, id) => [{ type: "Transfer", id }],
    }),
    getPublicTransfer: build.query<PublicTransfer, string>({
      async queryFn(slug) {
        await delay(150);
        const item = mockPublicTransfers[slug] ?? mockPublicTransfers["brand-2025"];
        if (slug === "missing") {
          return { error: { status: 404, data: "not found" } };
        }
        return { data: item };
      },
    }),
  }),
});

function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export const {
  useListTransfersQuery,
  useGetTransferQuery,
  useGetPublicTransferQuery,
} = transfersApi;
