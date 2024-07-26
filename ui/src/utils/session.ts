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

export function getSessionInfos() {
  const token = cookieStorage.getItem("karaberus_session");
  if (!token) {
    return null;
  }
  const payload = decodeJwt<KaraberusTokenPayload>(token);
  if (payload.exp < Date.now() / 1000) {
    return null;
  }
  return payload;
}
