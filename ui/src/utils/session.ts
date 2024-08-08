import { cookieStorage } from "@solid-primitives/storage";
import { decodeJwt } from "jose";
import { isTauri } from "./tauri";

export interface KaraberusJwtPayload {
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
  kara_ro: boolean;
  user: boolean;
  is_admin: boolean;
}

export const SESSION_TOKEN_NAME = "karaberus_session";

export function getSessionToken() {
  return (isTauri ? localStorage : cookieStorage).getItem(SESSION_TOKEN_NAME);
}

export function setSessionToken(token: string) {
  (isTauri ? localStorage : cookieStorage).setItem(SESSION_TOKEN_NAME, token);
}

export function removeSessionToken() {
  (isTauri ? localStorage : cookieStorage).removeItem(SESSION_TOKEN_NAME);
}

export function getSessionInfos() {
  const token = getSessionToken();
  if (!token) {
    return null;
  }
  const payload = decodeJwt<KaraberusJwtPayload>(token);
  if (payload.exp < Date.now() / 1000) {
    return null;
  }
  return payload;
}

export const isAdmin = () => getSessionInfos()?.is_admin ?? false;
