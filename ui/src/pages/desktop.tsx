import { useNavigate } from "@solidjs/router";
import { createEffect } from "solid-js";
import { getSessionToken } from "../utils/session";

export default function Desktop() {
  const navigate = useNavigate();

  createEffect(() => {
    const token = getSessionToken();
    location.href = "karaberus://?token=" + token;
    navigate("/");
  });

  return null;
}
