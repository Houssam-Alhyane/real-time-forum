const app = document.getElementById('app');

// Handle routing based on the URL path when the page loads
document.addEventListener('DOMContentLoaded', () => {
  if (!app) {
    console.error('App container not found');
    return;
  }
  router();
});

// Handle browser back/forward navigation buttons cleanly
window.addEventListener('popstate', () => {
  router();
});

/**
 * Client-Side Router
 */
function router() {
  const path = window.location.pathname;

  const flash = localStorage.getItem('flash_message');
  if (flash) {
    displayMessage(flash, false);
    localStorage.removeItem('flash_message');
  }

  if (path === '/register') {
    renderRegister();
  } else if (path === '/login') {
    renderLogin();
  } else if (path === '/') {
    renderSPAHome();
  } else {
    app.innerHTML = `
      <div class="not-found">
        <h2>404 — Page Not Found</h2>
        <p onclick="navigateTo('/')">Go back home</p>
      </div>
    `;
  }
}

/**
 * Navigation Helper
 */
function navigateTo(path) {
  window.history.pushState({}, '', path);
  router();
}


function displayMessage(message, isError = true) {
  const existingMessage = document.getElementById('form-message');
  if (existingMessage) existingMessage.remove();

  const messageDiv = document.createElement('div');
  messageDiv.id = 'form-message';
  messageDiv.classList.add(isError ? 'is-error' : 'is-success');
  messageDiv.innerText = message;

  const heading = document.querySelector('#app h2');
  if (heading) {
    //add messageDiv after h2
    heading.insertAdjacentElement('afterend', messageDiv);
  }
}

// ---------------- LOGIN UI ----------------
function renderLogin() {
  if (
    window.location.pathname !== '/login' &&
    window.location.pathname !== '/'
  ) {
    window.history.pushState({}, '', '/login');
  }

  app.innerHTML = `
    <div class="auth-shell">
      <h2>Login</h2>
      <input id="login-id" placeholder="Email or Nickname">
      <input id="login-pass" type="password" placeholder="Password">
      <button type="button" onclick="login()">Login</button>
      <p class="auth-switch" onclick="navigateTo('/register')">Create account</p>
    </div>
  `;
}

// ---------------- REGISTER UI ----------------
function renderRegister() {
  if (window.location.pathname !== '/register') {
    window.history.pushState({}, '', '/register');
  }

  app.innerHTML = `
    <div class="auth-shell">
      <h2>Register</h2>
      <input id="nickname" placeholder="nickname">
      <input id="first_name" placeholder="first name">
      <input id="last_name" placeholder="last name">
      <input id="age" type="number" placeholder="age">
      <select id="gender">
        <option value="">Select gender</option>
        <option value="male">male</option>
        <option value="female">female</option>
      </select>
      <input id="email" placeholder="email">
      <input id="password" type="password" placeholder="password">
      <input id="confirm_password" type="password" placeholder="confirm password">
      <button type="button" onclick="register()">Register</button>
      <p class="auth-switch" onclick="navigateTo('/login')">Login</p>
    </div>
  `;
}

// ---------------- SPA HOME DASHBOARD UI ----------------
function renderSPAHome() {
  if (window.location.pathname !== '/') {
    window.history.pushState({}, '', '/');
  }

  const isLoggedIn = document.cookie.includes('logged_in');

  app.innerHTML = `
    <div class="zone-home">
      <section id="posts-section">
        <h2>Forum Feed</h2>
        <div id="posts-container">Loading interactive feed...</div>
      </section>

      ${
        isLoggedIn
          ? `
        <section id="chat-sidebar">
          <h2>Active Chat</h2>
          <div id="users-list">Loading chat list...</div>
          <div id="private-messages-box">Select a user to view chat history</div>
          <div class="chat-input-row">
            <input id="msg-input" placeholder="Type a message...">
            <button onclick="sendPrivateMessage()">Send</button>
          </div>
          <p class="logout-link" onclick="handleLogout()">Logout from all pages</p>
        </section>
      `
          : `
        <section id="chat-sidebar">
          <h2>Join the conversation</h2>
          <p style="color: var(--text-muted); font-size: 0.9rem; margin-bottom: 18px;">
            Login or create an account to chat and post.
          </p>
          <button type="button" onclick="navigateTo('/login')">Login</button>
          <p class="auth-switch" onclick="navigateTo('/register')">Create account</p>
        </section>
      `
      }
    </div>
  `;

  if (typeof loadPosts === 'function') {
    loadPosts();
  }

  if (isLoggedIn && typeof initWebSocketsConnection === 'function') {
    initWebSocketsConnection();
  }
}

// ---------------- REGISTER ACTION ----------------
async function register() {
  const required = [
    'nickname',
    'first_name',
    'last_name',
    'age',
    'gender',
    'email',
    'password',
    'confirm_password',
  ];
  const data = {};

  for (let id of required) {
    const el = document.getElementById(id);
    if (!el || el.value.trim() === '') {
      displayMessage('All fields are required', true);
      return;
    }
    data[id] = el.value;
  }

  try {
    const res = await fetch('/register', {
      method: 'POST',
      body: new URLSearchParams(data),
    });

    const result = await res.json();

    if (!res.ok) {
      displayMessage(result.error || 'Registration failed', true);
      return;
    }

    localStorage.setItem(
      'flash_message',
      result.message || 'Account created successfully!'
    );
    navigateTo('/login');
  } catch (err) {
    console.error('Register network error:', err);
    displayMessage('Network error, please try again', true);
  }
}

// ---------------- LOGIN ACTION ----------------
async function login() {
  const loginInput = document.getElementById('login-id');
  const passwordInput = document.getElementById('login-pass');

  if (
    !loginInput ||
    !passwordInput ||
    loginInput.value.trim() === '' ||
    passwordInput.value.trim() === ''
  ) {
    displayMessage('Email/username and password are required', true);
    return;
  }

  const data = {
    login: loginInput.value,
    password: passwordInput.value,
  };

  try {
    const res = await fetch('/login', {
      method: 'POST',
      body: new URLSearchParams(data),
    });

    const result = await res.json();

    if (!res.ok) {
      displayMessage(result.error || 'Login failed', true);
      return;
    }

    localStorage.setItem('flash_message', result.message || 'Login successful');
    navigateTo('/');
  } catch (err) {
    console.error('Login network error:', err);
    displayMessage('Network error, please try again', true);
  }
}

// ---------------- LOGOUT ACTION ----------------
async function handleLogout() {
  try {
    await fetch('/logout', { method: 'POST' });
  } catch (err) {
    console.error('Logout network error:', err);
  }
  document.cookie =
    'logged_in=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
  navigateTo('/');
}
