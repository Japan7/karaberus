import type { Component } from 'solid-js';

import { Route, Router } from '@solidjs/router'
import Home from './Home'
import { login_path, oidc_callback_path } from './oidc'
import Login from './Login';
import OIDCCallback from './OIDCCallback';

const App: Component = () => {
  return (
    <Router>
      <Route path={login_path} component={Login} />
      <Route path={oidc_callback_path} component={OIDCCallback} />
      <Route path="/" component={Home} />
    </Router>
  );
};

export default App;
