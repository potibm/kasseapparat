import { z } from "zod";
import {
  AuthUser as AuthUserType,
  Session as SessionType,
} from "../types/auth.types";
import { UserDataSchema } from "@core/api/auth.schemas";
import { createLogger } from "@core/logger/logger";

const log = createLogger('Auth');

const LS_PREFIX = "kasseapparat.auth.";
export const AUTH_KEYS = {
  TOKEN: `${LS_PREFIX}token`,
  EXPIRY: `${LS_PREFIX}expiryDate`,
  USER: `${LS_PREFIX}userdata`,
};

const LSSessionSchema = z.object({
  token: z.jwt(),
  expiryDate: z.iso.datetime().transform((val) => new Date(val)),
});

export const getInitialSession = (): SessionType => {
  try {
    const raw = {
      token: localStorage.getItem(AUTH_KEYS.TOKEN),
      expiryDate: localStorage.getItem(AUTH_KEYS.EXPIRY),
    };

    const result = LSSessionSchema.safeParse(raw);

    if (!result.success) {
      return { token: null, expiryDate: null };
    }

    log.debug("LocalStorage Session restored", result.data);
    return result.data;
  } catch {
    return { token: null, expiryDate: null };
  }
};

export const getInitialUser = (): AuthUserType | null => {
  try {
    const data = localStorage.getItem(AUTH_KEYS.USER);
    if (!data) return null;

    const parsed = JSON.parse(data);

    const result = UserDataSchema.safeParse(parsed);

    if (!result.success) {
      log.warn("LocalStorage Userdata invalid. Clearing...");
      return null;
    }

    log.debug("LocalStorage Userdata restored", result.data);
    return result.data;
  } catch {
    return null;
  }
};

export const storeSession = (token: string, expiryDate: Date) => {
  localStorage.setItem(AUTH_KEYS.TOKEN, token);
  localStorage.setItem(AUTH_KEYS.EXPIRY, expiryDate.toISOString());
};

export const storeUser = (userdata: AuthUserType) => {
  localStorage.setItem(AUTH_KEYS.USER, JSON.stringify(userdata));
};

export const clearAuthStorage = () => {
  Object.values(AUTH_KEYS).forEach((key) => localStorage.removeItem(key));
};
