import { state } from './state.js';
import { navigateTo } from './routeer.js';
import { displayMessage } from './toast.js';

export async function handleLogout() {
  try {
    await fetch('/logout', { method: 'POST' });
  } catch (err) {
    console.error('Logout error:', err);
  }
  // Reset local state — server already cleared the session cookie
  state.auth = { authenticated: false, user: null };
  displayMessage('logout successfully', false);
  navigateTo('/');
}
