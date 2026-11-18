import type { User } from "@/lib/types/auth";
import type {
  AdminSetTransferStatusRequest,
  AdminSetUserActiveRequest,
  AdminTransferListItem,
} from "@/lib/types/admin";
import type { Transfer } from "@/lib/types/transfer";
import { baseApi } from "@/store/api/baseApi";

export const adminApi = baseApi.injectEndpoints({
  endpoints: (build) => ({
    listAdminUsers: build.query<User[], void>({
      query: () => "/api/v1/admin/users",
      providesTags: [{ type: "AdminUsers", id: "LIST" }],
    }),
    setUserActive: build.mutation<
      User,
      { id: string; body: AdminSetUserActiveRequest }
    >({
      query: ({ id, body }) => ({
        url: `/api/v1/admin/users/${id}`,
        method: "PATCH",
        body,
      }),
      invalidatesTags: [{ type: "AdminUsers", id: "LIST" }],
    }),
    listAdminTransfers: build.query<AdminTransferListItem[], void>({
      query: () => "/api/v1/admin/transfers",
      providesTags: [{ type: "AdminTransfers", id: "LIST" }],
    }),
    setTransferStatus: build.mutation<
      Transfer,
      { id: string; body: AdminSetTransferStatusRequest }
    >({
      query: ({ id, body }) => ({
        url: `/api/v1/admin/transfers/${id}`,
        method: "PATCH",
        body,
      }),
      invalidatesTags: [
        { type: "AdminTransfers", id: "LIST" },
        { type: "Transfers", id: "LIST" },
      ],
    }),
  }),
});

export const {
  useListAdminUsersQuery,
  useSetUserActiveMutation,
  useListAdminTransfersQuery,
  useSetTransferStatusMutation,
} = adminApi;
