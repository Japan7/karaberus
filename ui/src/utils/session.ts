import { cookieStorage } from "@solid-primitives/storage";
import { decodeJwt } from "jose";
import { IS_TAURI_DIST_BUILD } from "./tauri";

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

const storage = IS_TAURI_DIST_BUILD ? localStorage : cookieStorage;

export function getSessionToken() {
  return storage.getItem(SESSION_TOKEN_NAME);
}

export function setSessionToken(token: string) {
  storage.setItem(SESSION_TOKEN_NAME, token);
}

export function removeSessionToken() {
  storage.removeItem(SESSION_TOKEN_NAME);
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
