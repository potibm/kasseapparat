import { z } from "zod";

const UserRoleSchema = z.enum(["user", "admin"]);

export const LoginResponseSchema = z.object({
  access_token: z.jwt(),
  token_type: z.string(),
  expires_in: z.number(),
  role: UserRoleSchema,
  username: z.string(),
  gravatarUrl: z.string(),
  id: z.number(),
});

export const UserDataSchema = z.object({
  role: UserRoleSchema,
  username: z.string(),
  gravatarUrl: z.string(),
  id: z.number(),
});

export const RefreshTokenResponseSchema = z.object({
  access_token: z.jwt(),
  refresh_token: z.string(),
  token_type: z.string(),
  expires_in: z.number(),
});

export const StringResponseSchema = z.string();

export const SimpleResponseSchema = z.object({
  message: z.string().optional(),
  status: z.string().optional(),
  code: z.number().optional(),
});

export type LoginResponse = z.infer<typeof LoginResponseSchema>;
export type RefreshTokenResponse = z.infer<typeof RefreshTokenResponseSchema>;
export type SimpleResponse = z.infer<typeof SimpleResponseSchema>;
export type StringResponse = z.infer<typeof StringResponseSchema>;
