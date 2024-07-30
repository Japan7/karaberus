import { Route, Router } from "@solidjs/router";
import Layout from "./layout/Layout";
import FontsBrowse from "./routes/fonts/Browse";
import Home from "./routes/Home";
import KaraokeBrowse from "./routes/karaoke/Browse";
import KaraokeCreate from "./routes/karaoke/Create";
import KaraokeIssues from "./routes/karaoke/Issues";
import TagsArtist from "./routes/tags/Artist";
import TagsAuthor from "./routes/tags/Author";
import TagsMedia from "./routes/tags/Media";
import routes from "./utils/routes";

export default function App() {
  return (
    <Router root={Layout}>
      <Route path={routes.HOME} component={Home} />

      <Route path={routes.KARAOKE_CREATE} component={KaraokeCreate} />
      <Route path={routes.KARAOKE_BROWSE} component={KaraokeBrowse} />
      <Route path={routes.KARAOKE_ISSUES} component={KaraokeIssues} />

      <Route path={routes.TAGS_MEDIA} component={TagsMedia} />
      <Route path={routes.TAGS_ARTIST} component={TagsArtist} />
      <Route path={routes.TAGS_AUTHOR} component={TagsAuthor} />

      <Route path={routes.FONTS_BROWSE} component={FontsBrowse} />
    </Router>
  );
}
