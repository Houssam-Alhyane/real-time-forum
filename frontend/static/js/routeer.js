import { renderHome } from './pages/home.js';
import { renderLogin } from './pages/login.js';
import { renderRegister } from './pages/register.js';
import { renderError } from './pages/error.js';
import { state } from './state.js';

export function router() {
  const path = location.pathname;

  const publicPaths = ['/login', '/register'];
  const knownPaths = ['/', ...publicPaths];
  const auth = state.auth.authenticated;

  if (!knownPaths.includes(path)) {
    renderError(404);
    return;
  }

  if (!auth && !publicPaths.includes(path)) {
    navigateTo('/login');
    return;
  }

  if (path === '/') {
    renderHome();
  } else if (path === '/login') {
    renderLogin();
  } else if (path === '/register') {
    renderRegister();
  }
}

export function navigateTo(path) {
  history.pushState({}, '', path);
  router();
}
