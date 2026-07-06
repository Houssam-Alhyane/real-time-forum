export const state = {
  auth: {
    authenticated: false,
    id: null,
    nickname: null,
  },
};

export function resetAuth() {
  state.auth = {
    authenticated: false,
    id: null,
    nickname: null,
  };
}

export async function initAuth() {
  try {
    const res = await fetch('/api/me');
    if (!res.ok) throw new Error('auth check failed');
    state.auth = await res.json();
  } catch (err) {
    console.error('initAuth error:', err);
    resetAuth();
  }
}
