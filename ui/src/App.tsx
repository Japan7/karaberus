import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import { Route, Router } from "@solidjs/router";
import { CssBaseline, ThemeProvider, createTheme } from "@suid/material";
import Layout from "./layout/Layout";
import Home from "./routes/Home";
import OIDCCallback from "./routes/OIDCCallback";
import OIDCLogin from "./routes/OIDCLogin";
import routes from "./utils/routes";

const theme = createTheme({
  palette: {
    primary: { main: "#9966cc" },
    secondary: { main: "#e968a8" },
  },
});

export default function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Layout>
        <Router>
          <Route path={routes.HOME} component={Home} />
          <Route path={routes.OIDC_CALLBACK} component={OIDCCallback} />
          <Route path={routes.OIDC_LOGIN} component={OIDCLogin} />
        </Router>
      </Layout>
    </ThemeProvider>
  );
}
