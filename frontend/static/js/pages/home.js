import { renderNavbar } from '../navbar.js';
import { loadPosts } from '../posts.js';
import { renderChatUsers, openChatPanel } from '../Chatui.js';
import { initWebSocket, chatState } from '../ChatData.js';
import { filterSidebarHTML } from '../utils.js';

async function fetchChatUsers() {
  try {
    const response = await fetch('/api/users');
    if (!response.ok) {
      throw new Error('Failed to fetch chat users');
    }
    const users = await response.json();
    return users;
  } catch (error) {
    console.error('Failed to fetch chat users:', error);
    return [];
  }
}

export function renderHome(initialCategories = []) {
  const app = document.getElementById('app');
  if (location.pathname !== '/') history.pushState({}, '', '/');

  app.innerHTML = `
    ${renderNavbar()}
    <div class="container">
      <div class="sidebar-col">
        ${filterSidebarHTML()}
      </div>

      <main class="content">
        <div id="posts-container">Loading feed...</div>
      </main>

    </div>`;

  // Pre-apply filters selected elsewhere (e.g. the post page sidebar) so the
  // feed renders filtered on the first fetch.
  if (initialCategories.length > 0) {
    document.querySelectorAll('.sidebar input[type="checkbox"]').forEach((cb) => {
      cb.checked = initialCategories.includes(cb.value);
    });
  }

  // Initialize WebSocket connection
  initWebSocket();

  // Load initial chat users (will be updated via WebSocket)
  fetchChatUsers().then((users) => {
    renderChatUsers(users);
    if (chatState.activeUserId) {
      openChatPanel(chatState.activeUserId);
    }
  });

  loadPosts();
}
