import { baseApi } from "@/store/api/baseApi";
import type {
  CreateTransferRequest,
  PublicTransfer,
  Transfer,
  TransferDetail,
  TransferFile,
  TransferFilter,
  TransferListItem,
} from "@/lib/types/transfer";

export type { TransferFilter };

export const transfersApi = baseApi.injectEndpoints({
  endpoints: (build) => ({
    listTransfers: build.query<TransferListItem[], TransferFilter>({
      query: (status) => ({
        url: "/api/v1/transfers",
        params: { status },
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.map(({ id }) => ({ type: "Transfer" as const, id })),
              { type: "Transfers", id: "LIST" },
            ]
          : [{ type: "Transfers", id: "LIST" }],
    }),
    getTransfer: build.query<TransferDetail, string>({
      query: (id) => `/api/v1/transfers/${id}`,
      providesTags: (_r, _e, id) => [{ type: "Transfer", id }],
    }),
    getPublicTransfer: build.query<PublicTransfer, string>({
      query: (shortCode) => `/api/v1/public/transfers/${shortCode}`,
    }),
    createTransfer: build.mutation<Transfer, CreateTransferRequest>({
      query: (body) => ({
        url: "/api/v1/transfers",
        method: "POST",
        body,
      }),
      invalidatesTags: [{ type: "Transfers", id: "LIST" }],
    }),
    uploadTransferFile: build.mutation<
      TransferFile,
      { transferId: string; file: File }
    >({
      query: ({ transferId, file }) => {
        const body = new FormData();
        body.append("file", file);
        return {
          url: `/api/v1/transfers/${transferId}/files`,
          method: "POST",
          body,
        };
      },
      invalidatesTags: (_r, _e, arg) => [
        { type: "Transfer", id: arg.transferId },
        { type: "Transfers", id: "LIST" },
      ],
    }),
  }),
});

export const {
  useListTransfersQuery,
  useGetTransferQuery,
  useGetPublicTransferQuery,
  useCreateTransferMutation,
  useUploadTransferFileMutation,
} = transfersApi;
