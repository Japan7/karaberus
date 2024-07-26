import { cookieStorage } from "@solid-primitives/storage";
import { decodeJwt } from "jose";

export interface KaraberusTokenPayload {
  sub: string;
  exp: number;
  iat: number;
  name: string;
  nickname: string;
  picture: string;
  locale: string;
  preferred_username: string;
  email: string;
  email_verified: boolean;
  kara: boolean;
  user: boolean;
}

export const SESSION_TOKEN_NAME = "karaberus_session";

export function getSessionInfos() {
  const token = cookieStorage.getItem(SESSION_TOKEN_NAME);
  if (!token) {
    return null;
  }
  const payload = decodeJwt<KaraberusTokenPayload>(token);
  if (payload.exp < Date.now() / 1000) {
    return null;
  }
  return payload;
}

export function clearSession() {
  cookieStorage.removeItem(SESSION_TOKEN_NAME);
}
