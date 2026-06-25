import { state } from './state.js';
import { escapeHTML } from './utils.js';

export function renderNavbar() {
  return `
    <header class="navbar">
      <div class="logo" onclick="window._nav('/')">01Forum</div>
      <div class="auth-buttons">
        ${
          state.auth.authenticated
            ? `<span class="nav-username">${escapeHTML(
                state.auth.user?.nickname || ''
              )}</span>
               <button class="btn logout" onclick="window._logout()">Logout</button>`
            : `<button class="btn login"    onclick="window._nav('/login')">Login</button>
               <button class="btn register" onclick="window._nav('/register')">Register</button>`
        }
      </div>
    </header>`;
}
