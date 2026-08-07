<div align="center">

# 💬 Real Time Forum

**A single-page forum application with registration, posts, comments, reactions, and real-time private messaging over WebSockets.**

![Go](https://img.shields.io/badge/Go-1.25-30363D?style=flat-square&labelColor=00ADD8&logo=go&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-embedded-30363D?style=flat-square&labelColor=003B57&logo=sqlite&logoColor=white)
![WebSockets](https://img.shields.io/badge/WebSockets-realtime-30363D?style=flat-square&labelColor=010101)
![Vanilla JS](https://img.shields.io/badge/Vanilla%20JS-no%20frameworks-30363D?style=flat-square&labelColor=F7DF1E&logo=javascript&logoColor=black)

</div>

---

## 📖 About

A forum application rewritten as a **single-page application (SPA)**: one HTML file, with every page transition (login, register, feed, post view) handled entirely in JavaScript on the client, talking to a Go backend over HTTP and WebSockets.

The backend is built with the Go standard library, SQLite for persistence, and WebSockets for real-time features. The frontend is written in **vanilla JavaScript** — no React, Vue, or any other framework — with hand-rolled routing, DOM rendering, state management, and a fully responsive stylesheet. Built as a learning project covering real-time communication, session management, and full-stack development from scratch.

## ✨ Features

### 🔐 Authentication & Sessions

- **Registration** with nickname, age, gender, first name, last name, email, and password — fully validated server-side (required fields, 18–200 age range, valid email format, password 6–21 characters, unique email & nickname).
- **Login with either nickname or email** + password.
- Passwords are hashed with **bcrypt** before storage.
- Sessions use **UUID v4 tokens** stored in HttpOnly, SameSite cookies with a **24-hour expiry**.
- Logging in again **force-logs-out all other open tabs** of that user in real time.
- Log out from anywhere in the app — the session is revoked and all the user's sockets are disconnected.

### 📝 Posts, Comments & Reactions

- Create posts tagged with **one or more categories** (General, Programming, Gaming, Movies, Sports, Anime).
- **Paginated feed** (newest first) with category filtering and cursor-based loading via the `X-Max-Post-Id` header.
- 👍/👎 **like & dislike reactions** on both posts and comments, with toggle behavior (re-act to switch, re-click to remove).
- Comments are loaded on demand per post, with their own pagination, author nickname, timestamps, and reaction counts.

### ⚡ Real-time Private Messaging

- A persistent sidebar lists all users with live **online/offline presence**, sorted by most recent message (Discord-style); users with no history are sorted alphabetically.
- **Paginated chat history** — the last 10 messages load first; scrolling to the top loads 10 more. Scroll handling is **throttled** so it never floods the backend, and pagination is anchored on the oldest loaded message so live updates never cause duplicates.
- **Typing indicators** broadcast to the receiver while you type.
- New messages and presence changes are pushed instantly over WebSockets — no page refresh.
- **Multi-tab support**: a user can be connected from several tabs; presence flips to offline only when the *last* tab disconnects.

### 📱 Responsiveness & Pagination

Clean, zoom-safe layout from desktop down to small phones, with posts, comments, and chat history loading incrementally.

### 🛡️ Backend Protections

- **Rate limiting**: fixed-window per-IP limiter (60 requests / 3 seconds) applied to all routes.
- **Auth middleware** verifies the session on every protected route and updates the user's `last_seen` presence timestamp.
- Full input validation on all endpoints, and **parameterized SQL queries** throughout (no string-built queries).

## 🧰 Tech Stack

| Layer    | Technology |
| -------- | ---------- |
| Language | [Go](https://go.dev/) 1.25 — standard library `net/http` |
| WebSocket | [Gorilla WebSocket](https://github.com/gorilla/websocket) |
| Database | [SQLite](https://www.sqlite.org/) via [`mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3) |
| Password hashing | [`golang.org/x/crypto/bcrypt`](https://pkg.go.dev/golang.org/x/crypto/bcrypt) |
| Session tokens | [`gofrs/uuid`](https://github.com/gofrs/uuid) |
| Frontend | Vanilla JavaScript (ES modules, no frameworks) + HTML/CSS |

**No frontend frameworks or libraries** are used — routing, DOM manipulation, and state management are all hand-written.

## ⚙️ Architecture

### Backend

- Serves the SPA and static assets; unknown routes and `/api` misses return proper 404s.
- REST-style endpoints for auth, posts, comments, and reactions, all sharing JSON helpers.
- **WebSocket hub** ([`handlers/websocket.go`](backend/handlers/websocket.go)):
  - A registry `map[userID][]*wsClient` tracks every open connection, guarded by a **mutex** (multiple tabs per user supported).
  - Each connection runs two **goroutines** — a read pump and a write pump — with periodic **ping keepalives** (30s) to detect dead connections.
  - `sync.Once` guarantees each connection is cleaned up exactly once; online/offline events are broadcast when the first/last tab connects/disconnects.
  - Concurrency is handled with goroutines and mutexes.

### WebSocket protocol

JSON envelopes with a `type` field over `/ws`.

| Direction | Type | Purpose |
| --------- | ---- | ------- |
| Client → Server | `get_history` | Request a page of chat history (`partner_id`, `offset`) |
| Client → Server | `send_message` | Deliver a private message (`sender_id`, `receiver_id`, `content`) |
| Client → Server | `typing` | Notify the receiver that you are typing (`receiver_id`) |
| Server → Client | `history_response` | A page of messages with `has_more` flag |
| Server → Client | `new_message` | A new private message, delivered live |
| Server → Client | `user_online` / `user_offline` | Presence changes (`user_id`, `nickname`) |
| Server → Client | `typing` | The peer is typing (`sender_id`, `sender_nickname`) |
| Server → Client | `force_logout` | Session superseded — log out all other tabs |
| Server → Client | `error` | Socket-level error message |

### Frontend

- `router.js` — client-side router that swaps views without page loads.
- `ChatData.js` — owns the WebSocket lifecycle, chat/user state, and all socket message handling.
- `Chatui.js` — chat DOM rendering (sidebar, message panel, bubbles).
- `app-events.js` — centralized event delegation via `data-action` attributes.
- `auth.js` — auth state sync and force-logout handling; `state.js` — shared app state; `toast.js` — user notifications; `typing.js` — typing indicator UI.

## 🗄️ Database

SQLite database (auto-created at `backend/database/forum.db`, git-ignored). The schema is defined in [`backend/database/schema.sql`](backend/database/schema.sql) and applied automatically on startup.

| Table | Purpose |
| ----- | ------- |
| `users` | Accounts (nickname, profile fields, bcrypt hash, `last_seen` presence) |
| `sessions` | Session tokens with expiry timestamps |
| `posts` | Forum posts |
| `categories` | Six seeded categories |
| `post_categories` | Many-to-many posts ↔ categories |
| `post_reactions` | Likes/dislikes on posts |
| `comments` | Comments on posts |
| `comment_reactions` | Likes/dislikes on comments |
| `messages` | Private chat messages |

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) `1.25` or newer
- A C compiler (e.g. GCC) — required by `mattn/go-sqlite3`

### Run

Start the server from the **project root** (database paths are relative):

```bash
go run main.go
```

Then open [http://localhost:8082](http://localhost:8082) in your browser.

### Run with Docker

Requires [Docker](https://www.docker.com/) with Compose.

```bash
docker compose up -d
```

Then open [http://localhost:9070](http://localhost:9070) in your browser (mapped from the container's port `8082`). No local Go toolchain or C compiler is needed — the app builds inside the container.

- The SQLite database persists in the `forum-data` volume across restarts.
- Stop with `docker compose down`; stop **and reset** the database with `docker compose down -v`.

## 🗂️ Project Structure

```
├── backend
│   ├── database       # SQLite schema, init, and query helpers
│   ├── handlers       # HTTP + WebSocket route handlers
│   ├── middleware     # Auth & rate limiting
│   ├── routing        # Route registration
│   └── types          # Shared Go types/structs
├── frontend
│   ├── index.html     # The single HTML entry point
│   └── static
│       ├── css        # Stylesheets
│       └── js
│           ├── pages/         # Per-view render functions (login, register, home, post, error)
│           ├── ChatData.js    # Chat state, normalization, WebSocket lifecycle
│           ├── Chatui.js      # Chat DOM rendering (sidebar, panel, messages)
│           ├── app-events.js  # Delegated event wiring (clicks, keydown, socket events)
│           ├── router.js      # Client-side router (SPA navigation)
│           ├── auth.js        # Auth state sync, logout, force-logout handling
│           ├── posts.js / comments.js / reactions.js
│           ├── state.js       # Shared client-side app state
│           └── main.js        # App bootstrap
├── go.mod / go.sum
└── main.go
```

## 🎓 Learning Objectives

This project covers:

- HTML, HTTP, sessions & cookies, CSS & responsive design
- Backend/Frontend separation and the DOM
- Go **goroutines and mutexes** for concurrency
- WebSockets, both server-side (Go) and client-side (JS)
- Real-time protocol design (typed JSON events, presence, pagination)
- SQL and database manipulation
- Building a single-page application without a frontend framework

## 👥 Authors

- houssam alhyane — [GitHub](https://github.com/Houssam-Alhyane)
- elmehdi rezoug — [GitHub](https://github.com/elmehdi-rezoug)
