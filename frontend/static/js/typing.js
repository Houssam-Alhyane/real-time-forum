import { chatState, parseUserId, typingState } from './ChatData.js';

const TYPING_STOP_DELAY_MS = 500;
//send data to websocket
function sendTypingSignal(type, receiverId) {
  if (!chatState.socket || chatState.socket.readyState !== WebSocket.OPEN)
    return;
  chatState.socket.send(JSON.stringify({ type, receiver_id: receiverId }));
}

//check
export function notifyTyping() {
  const partnerId = chatState.activeUserId;
  if (!partnerId) return;
  const now = Date.now();
  sendTypingSignal('typing_start', partnerId);
  typingState.lastStartSentAt = now;
  if (typingState.stopTimeoutId) {
    clearTimeout(typingState.stopTimeoutId);
  }
  typingState.stopTimeoutId = setTimeout(() => {
    sendTypingSignal('typing_stop', partnerId);
    typingState.lastStartSentAt = 0;
  }, TYPING_STOP_DELAY_MS);
}

//check if user close chat or go to another new tab
export function stopTypingNow() {
  //receiver
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
//render typing to user
export function showTypingIndicator(nickname) {
  const el = getTypingIndicatorEl();
  if (!el) return;
  const nameEl = el.querySelector('.chat-typing-name');
  if (nameEl) nameEl.textContent = `${nickname} is typing`;
  el.classList.add('visible');
}

//remove typing to user
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
//reset typing
export function resetTypingUI() {
  stopTypingNow();
  hideTypingIndicator();
  if (typingState.remoteTypingTimeoutId) {
    clearTimeout(typingState.remoteTypingTimeoutId);
    typingState.remoteTypingTimeoutId = null;
  }
}
