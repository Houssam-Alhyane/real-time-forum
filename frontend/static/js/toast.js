export function displayMessage(message, isError = true) {
  const existing = document.getElementById('toast-notification');
  if (existing) existing.remove();

  const toast = document.createElement('div');
  toast.id = 'toast-notification';
  toast.classList.add('toast', isError ? 'toast-error' : 'toast-success');
  toast.innerHTML = `
    <span class="toast-icon">${isError ? '✕' : '✓'}</span>
    <span class="toast-text">${message}</span>
  `;

  document.body.appendChild(toast);
  requestAnimationFrame(() => toast.classList.add('toast-visible'));

  setTimeout(() => {
    toast.classList.remove('toast-visible');
    toast.addEventListener('transitionend', () => toast.remove(), {
      once: true,
    });
  }, 3000);
}
