

let allPosts = []; // cache for client-side filtering

// ---------------- LOAD POSTS ----------------
async function loadPosts() {
  const container = document.getElementById('posts-container');
  if (!container) return;

  try {
    const res = await fetch('/api/posts');
    if (!res.ok) throw new Error('Failed to fetch posts');
    allPosts = await res.json();
    renderPosts(allPosts);
  } catch (err) {
    console.error('loadPosts error:', err);
    container.innerHTML = `<p class="error-text">Failed to load posts. Please try again.</p>`;
  }
}

// ---------------- FILTER ----------------
function filterPosts() {
  const checked = Array.from(
    document.querySelectorAll('.sidebar input[type=checkbox]:checked')
  ).map((cb) => cb.value);

  if (checked.length === 0) {
    renderPosts(allPosts);
  } else {
    const filtered = allPosts.filter((p) => checked.includes(p.category_name));
    renderPosts(filtered);
  }
}

function clearFilters() {
  document
    .querySelectorAll('.sidebar input[type=checkbox]')
    .forEach((cb) => (cb.checked = false));
  renderPosts(allPosts);
}

// ---------------- RENDER POSTS ----------------
function renderPosts(posts) {
  const container = document.getElementById('posts-container');
  if (!container) return;

  const isLoggedIn = document.cookie.includes('logged_in');

  const createBtn = isLoggedIn
    ? `<button class="btn primary create-post-btn" onclick="renderCreatePostForm()">+ New Post</button>`
    : '';

  if (!posts || posts.length === 0) {
    container.innerHTML = `${createBtn}<p class="empty-feed">No posts found.</p>`;
    return;
  }

  const cards = posts
    .map(
      (p) => `
    <article class="post">
      <div class="post-header">
        <h3>${escapeHTML(p.title)}</h3>
        <div class="post-categories">
          <span class="category">${escapeHTML(p.category_name)}</span>
        </div>
      </div>
      <div class="post-body">
        <p>${escapeHTML(p.content)}</p>
      </div>
    </article>
  `
    )
    .join('');

  container.innerHTML = `${createBtn}<div class="posts-list">${cards}</div>`;
}

// ---------------- CREATE POST FORM ----------------
function renderCreatePostForm() {
  const container = document.getElementById('posts-container');
  if (!container) return;

  fetch('/api/categories')
    .then((res) => {
      if (!res.ok) throw new Error();
      return res.json();
    })
    .then((categories) => {
      const options = categories
        .map((c) => `<option value="${c.id}">${escapeHTML(c.name)}</option>`)
        .join('');

      container.innerHTML = `
        <div class="create-post-form">
          <h3>Create a New Post</h3>
          <input id="post-title" placeholder="Post title" maxlength="200">
          <select id="post-category">
            <option value="">Select a category</option>
            ${options}
          </select>
          <textarea id="post-content" placeholder="Write your post..." rows="5"></textarea>
          <div class="form-actions">
            <button type="button" class="btn primary" onclick="submitPost()">Publish</button>
            <button type="button" class="btn"         onclick="loadPosts()">Cancel</button>
          </div>
        </div>`;
    })
    .catch(() => {
      container.innerHTML = `<p class="error-text">Failed to load categories.</p>`;
    });
}

// ---------------- SUBMIT POST ----------------
async function submitPost() {
  const title = document.getElementById('post-title')?.value.trim();
  const content = document.getElementById('post-content')?.value.trim();
  const categoryID = document.getElementById('post-category')?.value;

  if (!title || !content || !categoryID) {
    displayMessage('All fields are required', true);
    return;
  }

  try {
    const res = await fetch('/api/posts/create', {
      method: 'POST',
      body: new URLSearchParams({ title, content, category_id: categoryID }),
    });
    const result = await res.json();

    if (!res.ok) {
      if (res.status === 401) {
        // Call logout to delete the session from DB, then clear cookie and redirect
        try {
          await fetch('/logout', { method: 'POST' });
        } catch (_) {}
        document.cookie =
          'logged_in=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
        localStorage.setItem(
          'flash_message',
          'Session expired. Please log in again.'
        );
        navigateTo('/login');
        return;
      }
      displayMessage(result.error || 'Failed to create post', true);
      return;
    }

    await loadPosts();
    displayMessage(result.message || 'Post created!', false);
  } catch (err) {
    console.error('submitPost error:', err);
    displayMessage('Network error. Please try again.', true);
  }
}

// ---------------- UTILITY ----------------
function escapeHTML(str) {
  if (typeof str !== 'string') return '';
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}
