import { chatState, parseUserId, typingState } from './ChatData.js';

//send status to websocket
function sendTypingSignal(receiverId) {
  if (!chatState.socket || chatState.socket.readyState !== WebSocket.OPEN)
    return;
  chatState.socket.send(
    JSON.stringify({ type: 'typing', receiver_id: receiverId })
  );
}
//check if the user typing
export function notifyTyping() {
  const partnerId = chatState.activeUserId;
  if (!partnerId) return;
  sendTypingSignal(partnerId);
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
}

//reset typing
export function resetTypingUI() {
  hideTypingIndicator();
  if (typingState.remoteTypingTimeoutId) {
    clearTimeout(typingState.remoteTypingTimeoutId);
    typingState.remoteTypingTimeoutId = null;
  }
}
