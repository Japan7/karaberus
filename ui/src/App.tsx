import { Route, Router } from "@solidjs/router";
import Home from "./routes/Home";
import Login from "./routes/Login";
import OIDCCallback from "./routes/OIDCCallback";
import { LOGIN_PATH, OIDC_CALLBACK_PATH } from "./utils/oidc";

export default function App() {
  return (
    <Router>
      <Route path={LOGIN_PATH} component={Login} />
      <Route path={OIDC_CALLBACK_PATH} component={OIDCCallback} />
      <Route path="/" component={Home} />
    </Router>
  );
}
