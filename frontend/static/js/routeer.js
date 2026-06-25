import { renderHome } from './pages/home.js';
import { renderLogin } from './pages/login.js';
import { renderRegister } from './pages/register.js';
import { renderError } from './pages/error.js';

export function router() {
  const path = window.location.pathname;

  if (path === '/') {
    renderHome();
  } else if (path === '/login') {
    renderLogin();
  } else if (path === '/register') {
    renderRegister();
  } else {
    renderError(404);
  }
}

export function navigateTo(path) {
  window.history.pushState({}, '', path);
  router();
}
