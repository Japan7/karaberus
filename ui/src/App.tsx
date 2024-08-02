import { Route, Router } from "@solidjs/router";
import Layout from "./layout/Layout";
import FontsBrowse from "./routes/fonts/Browse";
import Home from "./routes/Home";
import KaraokeBrowseId from "./routes/karaoke/browse/Id";
import KaraokeBrowse from "./routes/karaoke/browse/Index";
import KaraokeIssues from "./routes/karaoke/Issues";
import KaraokeNew from "./routes/karaoke/New";
import TagsArtist from "./routes/tags/Artist";
import TagsAuthor from "./routes/tags/Author";
import TagsMedia from "./routes/tags/Media";
import routes from "./utils/routes";

export default function App() {
  return (
    <Router root={Layout}>
      <Route path={routes.HOME} component={Home} />

      <Route path={routes.KARAOKE_NEW} component={KaraokeNew} />
      <Route path={routes.KARAOKE_BROWSE}>
        <Route path="/" component={KaraokeBrowse} />
        <Route path="/:id" component={KaraokeBrowseId} />
      </Route>
      <Route path={routes.KARAOKE_ISSUES} component={KaraokeIssues} />

      <Route path={routes.TAGS_MEDIA} component={TagsMedia} />
      <Route path={routes.TAGS_ARTIST} component={TagsArtist} />
      <Route path={routes.TAGS_AUTHOR} component={TagsAuthor} />

      <Route path={routes.FONTS_BROWSE} component={FontsBrowse} />
    </Router>
  );
}
