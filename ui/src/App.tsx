import { Route, Router } from "@solidjs/router";
import Layout from "./layout/Layout";
import Home from "./routes/Home";
import routes from "./utils/routes";

export default function App() {
  return (
    <Layout>
      <Router>
        <Route path={routes.HOME} component={Home} />
      </Router>
    </Layout>
  );
}
