import { state, initAuth } from '../state.js';
import { navigateTo } from '../router.js';
import { renderNavbar } from '../navbar.js';
import { displayMessage } from '../toast.js';

// Demo accounts seeded by backend/database/seed.go — keep in sync with it.
const DEMO_PASSWORD = 'demo1234';
const DEMO_USERS = [
  { nickname: 'Sara', name: 'Sara Johnson', initials: 'S', colors: ['#6e6bf0', '#8b5cf6'] },
  { nickname: 'Ahmed', name: 'Ahmed Ali', initials: 'A', colors: ['#f59e0b', '#ef4444'] },
  { nickname: 'Nicole', name: 'Nicole Martin', initials: 'N', colors: ['#ec4899', '#8b5cf6'] },
  { nickname: 'Alex', name: 'Alex Brown', initials: 'A', colors: ['#3ecf8e', '#14b8a6'] },
];

export function renderLogin() {
  const app = document.getElementById('app');

  if (state.auth.authenticated) {
    navigateTo('/');
    return;
  }
  if (location.pathname !== '/login') history.pushState({}, '', '/login');

  app.innerHTML = `
    ${renderNavbar()}
    <div class="auth-shell">
      <h2>Login</h2>
      <input id="login-id" maxlength="50"  minlength="2"   placeholder="Email or Nickname">
      <input id="login-pass" type="password"  maxlength="21"  minlength="6"  placeholder="Password">
      <button type="button" data-action="login">Login</button>

      <div class="demo-divider"><span>or try a demo account</span></div>
      <div class="demo-users">
        ${DEMO_USERS.map((u) => `
          <button type="button" class="demo-user" data-action="demo-login" data-nickname="${u.nickname}" title="Log in as ${u.name}">
            <span class="demo-avatar" style="--avatar-a:${u.colors[0]};--avatar-b:${u.colors[1]}">${u.initials}</span>
            <span class="demo-meta">
              <strong>${u.nickname}</strong>
              <small>${u.name}</small>
            </span>
          </button>`).join('')}
      </div>
    </div>`;
}

async function submitLogin(loginValue, password) {
  try {
    const res = await fetch('/login', {
      method: 'POST',
      body: new URLSearchParams({
        login: loginValue,
        password,
      }),
    });
    const result = await res.json();
    if (!res.ok) {
      displayMessage(result.error || 'Login failed', true);
      return false;
    }

    await initAuth();
    displayMessage('login successfully', false);
    navigateTo('/');
    return true;
  } catch (err) {
    console.error('Login error:', err);
    displayMessage('Network error, please try again', true);
    return false;
  }
}

export async function login() {
  const loginInput = document.getElementById('login-id');
  const passInput = document.getElementById('login-pass');

  if (!loginInput?.value.trim() || !passInput?.value.trim()) {
    displayMessage('Email/username and password are required', true);
    return;
  }

  await submitLogin(loginInput.value, passInput.value);
}

export async function demoLogin(nickname) {
  await submitLogin(nickname, DEMO_PASSWORD);
}
