import { Route, Router } from "@solidjs/router";
import Layout from "./layout/Layout";
import FontsBrowse from "./routes/fonts/Browse";
import FontsUpload from "./routes/fonts/Upload";
import Home from "./routes/Home";
import KaraokeBrowse from "./routes/karaoke/Browse";
import KaraokeIssues from "./routes/karaoke/Issues";
import KaraokeUpload from "./routes/karaoke/Upload";
import routes from "./utils/routes";

export default function App() {
  return (
    <Router root={Layout}>
      <Route path={routes.HOME} component={Home} />

      <Route path={routes.KARAOKE_UPLOAD} component={KaraokeUpload} />
      <Route path={routes.KARAOKE_BROWSE} component={KaraokeBrowse} />
      <Route path={routes.KARAOKE_ISSUES} component={KaraokeIssues} />

      <Route path={routes.FONTS_UPLOAD} component={FontsUpload} />
      <Route path={routes.FONTS_BROWSE} component={FontsBrowse} />
    </Router>
  );
}
