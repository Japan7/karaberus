import { useNavigate } from "@solidjs/router";
import { createEffect } from "solid-js";

export default function Redirect() {
  const navigate = useNavigate();

  createEffect(() => {
    const params = new URLSearchParams(location.search);
    const href = params.get("href");
    if (href) {
      setTimeout(() => navigate("/"));
      location.href = href;
    } else {
      navigate("/");
    }
  });

  return null;
}
