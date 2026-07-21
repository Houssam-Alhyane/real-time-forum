import { navigateTo } from './routeer.js';
import { renderError } from './pages/error.js';
import { displayMessage } from './toast.js';
import { postCardHTML, commentHTML, loadMoreCommentsBtnHTML } from './utils.js';

let offset = 0;
let hasMore = true;
const pageSize = 10;

function updateLoadMoreButton(postId) {
  const wrap = document.getElementById('comments-load-more-wrap');
  if (!wrap) return;
  wrap.outerHTML = loadMoreCommentsBtnHTML(hasMore, postId);
}

export async function fetchPostAndComments(postId) {
  offset = 0;
  hasMore = true;

  try {
    const [postResponse, commentsResponse] = await Promise.all([
      fetch(`/api/posts/${postId}`),
      fetch(
        `/api/posts/comments?post_id=${postId}&offset=${offset}&limit=${pageSize + 1}`,
      ),
    ]);

    if (!postResponse.ok || !commentsResponse.ok) {
      if (postResponse.status === 404) {
        renderError(404);
        return;
      }
      throw new Error('Failed to fetch post or comments');
    }

    const [post, commentsData] = await Promise.all([
      postResponse.json(),
      commentsResponse.json(),
    ]);
    const allComments = commentsData.comments || [];

    hasMore = allComments.length > pageSize;
    const comments = allComments.slice(0, pageSize);
    offset += comments.length;

    const postContainer = document.getElementById('post-container');
    if (postContainer) {
      postContainer.innerHTML = `<article class="post">${postCardHTML(post, {
        showCommentLink: false,
      })}</article>`;
    }

    renderComments(comments, false);
    updateLoadMoreButton(postId);
  } catch (error) {
    console.error('Error fetching post and comments:', error);
    displayMessage('Failed to load post and comments', true);
  }
}

async function fetchComments(postId) {
  if (!hasMore) return [];

  const response = await fetch(
    `/api/posts/comments?post_id=${postId}&offset=${offset}&limit=${pageSize + 1}`,
  );

  if (!response.ok) {
    throw new Error('Failed to fetch more comments');
  }

  const data = await response.json();
  const allComments = data.comments || [];

  hasMore = allComments.length > pageSize;
  const comments = allComments.slice(0, pageSize);
  offset += comments.length;

  return comments;
}

export async function loadMoreComments(postId) {
  try {
    const comments = await fetchComments(postId);
    if (comments.length > 0) {
      renderComments(comments, true);
    }
    updateLoadMoreButton(postId);
  } catch (error) {
    console.error('Error loading more comments:', error);
    displayMessage('Failed to load more comments', true);
  }
}

export function renderComments(comments, append = false) {
  const commentsContainer = document.getElementById('comments-container');
  if (!commentsContainer) return;

  if (!append) {
    commentsContainer.innerHTML = '';
  }

  if (comments.length === 0 && !append) {
    commentsContainer.innerHTML =
      '<p class="no-comments">No comments yet. Be the first to comment!</p>';
    return;
  }

  const commentsHTML = comments.map((comment) => commentHTML(comment)).join('');
  commentsContainer.insertAdjacentHTML('beforeend', commentsHTML);
}

export async function submitComment(postId) {
  const commentInputEl = document.getElementById('comment-input');
  const commentInput = commentInputEl?.value.trim() || '';

  if (!commentInput) {
    displayMessage('Comment cannot be empty', true);
    return;
  }

  const payload = new URLSearchParams({
    post_id: String(postId),
    content: commentInput,
  });

  try {
    const response = await fetch('/api/posts/comment', {
      method: 'POST',
      body: payload,
    });

    const data = await response.json();
    if (!response.ok) {
      displayMessage(data.message || 'Failed to submit comment', true);
      return;
    }

    if (commentInputEl) {
      commentInputEl.value = '';
    }

    displayMessage('Comment submitted successfully', false);
    navigateTo(`/posts/${postId}`);
  } catch (error) {
    displayMessage('Failed to submit comment', true);
  }
}
