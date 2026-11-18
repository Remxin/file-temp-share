export type User = {
  id: string;
  email: string;
  role: string;
  is_active: boolean;
  created_at: string;
};

export type AuthResponse = {
  token: string;
  expires_at: string;
  user: User;
};

export type MessageResponse = {
  message: string;
};

export type ApiErrorBody = {
  error: string;
  details?: { field: string; message: string }[];
  fields?: Record<string, string>;
};

export type RegisterRequest = {
  email: string;
  password: string;
};

export type LoginRequest = {
  email: string;
  password: string;
};

export type PasswordResetRequest = {
  email: string;
};

export type PasswordResetConfirmRequest = {
  token: string;
  new_password: string;
};

export type ChangePasswordRequest = {
  current_password: string;
  new_password: string;
};
