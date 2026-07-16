import { state } from '../state.js';
import { fetchPostAndComments } from '../comments.js';
import { renderNavbar } from '../navbar.js';

export function renderPostPage(postId) {
  // target the app container
  const app = document.getElementById('app');

  // inject the navbar and post container into the app
  app.innerHTML = `
    ${renderNavbar()}
    <div id="post-back-wrap">
      <button class="btn post-back-btn" data-action="nav" data-target="/">&larr; Back to Posts</button>
    </div>
    <div id="post-container"></div>
    <div id="comments-container"></div>
    <div id="comment-form-container">
      <textarea id="comment-input" placeholder="Write a comment..."></textarea>
      <button class="btn" data-action="submit-comment" data-post-id="${postId}">Submit Comment</button>
    </div>
    <div id="comments-load-more-wrap"></div>
  `;
  fetchPostAndComments(postId);
}
