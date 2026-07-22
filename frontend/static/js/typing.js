import { chatState, parseUserId, typingState, REMOTE_TYPING_TIMEOUT_MS } from './ChatData.js';

const TYPING_START_THROTTLE_MS = 2000;

const TYPING_STOP_DELAY_MS = 2000;


function sendTypingSignal(type, receiverId) {
  if (!chatState.socket || chatState.socket.readyState !== WebSocket.OPEN)
    return;
  chatState.socket.send(JSON.stringify({ type, receiver_id: receiverId }));
}

export function notifyTyping() {
  const partnerId = chatState.activeUserId;
  if (!partnerId) return;

  const now = Date.now();
  if (now - typingState.lastStartSentAt >= TYPING_START_THROTTLE_MS) {
    sendTypingSignal('typing_start', partnerId);
    typingState.lastStartSentAt = now;
  }

  if (typingState.stopTimeoutId) clearTimeout(typingState.stopTimeoutId);
  typingState.stopTimeoutId = setTimeout(() => {
    typingState.stopTimeoutId = null;
    sendTypingSignal('typing_stop', partnerId);
    typingState.lastStartSentAt = 0;
  }, TYPING_STOP_DELAY_MS);
}

export function stopTypingNow() {
  const partnerId = chatState.activeUserId;

  if (typingState.stopTimeoutId) {
    clearTimeout(typingState.stopTimeoutId);
    typingState.stopTimeoutId = null;
  }

  if (partnerId && typingState.lastStartSentAt !== 0) {
    sendTypingSignal('typing_stop', partnerId);
  }
  typingState.lastStartSentAt = 0;
}

function getTypingIndicatorEl() {
  return document.getElementById('chat-typing-indicator');
}

export function showTypingIndicator(nickname) {
  const el = getTypingIndicatorEl();
  if (!el) return;
  const nameEl = el.querySelector('.chat-typing-name');
  if (nameEl) nameEl.textContent = `${nickname} is typing`;
  el.classList.add('visible');
}

export function hideTypingIndicator() {
  const el = getTypingIndicatorEl();
  if (!el) return;
  el.classList.remove('visible');
}

export function attachTypingListeners() {
  const input = document.getElementById('chat-input');
  if (!input) return;

  input.addEventListener('input', notifyTyping);
  input.addEventListener('blur', stopTypingNow);
}

export function resetTypingUI() {
  stopTypingNow();
  hideTypingIndicator();
  if (typingState.remoteTypingTimeoutId) {
    clearTimeout(typingState.remoteTypingTimeoutId);
    typingState.remoteTypingTimeoutId = null;
  }
}
