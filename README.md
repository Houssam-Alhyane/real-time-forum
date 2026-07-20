# Real Time Forum

A single-page forum application with registration, posts, comments, and
real-time private messaging over WebSockets.

## Overview

This project builds on a previous forum implementation, rewritten as a
**single-page application (SPA)**: there is one HTML file, and every page
transition (login, register, feed, post view) is handled entirely in
JavaScript on the client, talking to a Go backend over HTTP and WebSockets.

## Tech Stack

| Layer    | Technology                                     |
| -------- | ---------------------------------------------- |
| Database | SQLite                                         |
| Backend  | Go (net/http, Gorilla WebSocket)               |
| Frontend | Vanilla JavaScript (ES modules, no frameworks) |
| Markup   | HTML (single page)                             |
| Styling  | CSS                                            |

No frontend frameworks or libraries (React, Vue, Angular, etc.) are used —
all DOM manipulation, routing, and state management are handled by hand in
plain JavaScript.

## Features

### Registration & Login

- Users register with: nickname, age, gender, first name, last name, email,
  and password.
- Users can log in with **either** their nickname or email, combined with
  their password.
- Sessions are managed via cookies; unauthenticated users only see the
  login/register screens.
- Users can log out from any page in the app.
- Passwords are hashed with `bcrypt` before being stored.

### Posts & Comments

- Authenticated users can create posts, each assigned to one or more
  categories.
- Posts are displayed in a feed.
- Comments are only loaded/shown when a specific post is opened.
- Users can comment on any post.

### Private Messages (real-time)

- A persistent sidebar lists all other users, showing **online/offline**
  status.
  - Sorted by most recent message (like Discord); users with no message
    history are sorted alphabetically.
- Clicking a user opens a chat panel and loads the most recent messages.
- **Chat history pagination**: the last 10 messages load initially; scrolling
  to the top loads 10 more. Scroll handling is **throttled**, not fired on
  every scroll event, to avoid flooding the backend with requests.
- Each message displays:
  - Timestamp
  - Sender's nickname
- New messages, and user online/offline status changes, are pushed to
  connected clients instantly via WebSockets — no page refresh required.

## Architecture

```
├── backend
│   ├── database      # SQLite schema, init, and query helpers
│   ├── handlers       # HTTP + WebSocket route handlers
│   ├── middleware      # Auth & rate limiting
│   ├── routing         # Route registration
│   └── types           # Shared Go types/structs
├── frontend
│   ├── index.html      # The single HTML entry point
│   └── static
│       ├── css          # Stylesheets
│       └── js
│           ├── pages        # Per-view render functions (login, register, home, post, error)
│           ├── ChatData.js  # Chat state, normalization helpers, WebSocket lifecycle
│           ├── Chatui.js    # Chat DOM rendering (sidebar, panel, messages)
│           ├── app-events.js # Delegated event wiring (clicks, keydown, socket events)
│           ├── routeer.js    # Client-side router (SPA navigation)
│           ├── auth.js       # Auth state sync, logout, force-logout handling
│           ├── posts.js / comments.js / reactions.js
│           ├── state.js      # Shared client-side app state
│           └── main.js       # App bootstrap
├── go.mod / go.sum
└── main.go
```

### Backend

- Serves the SPA and static assets.
- Exposes REST-style endpoints for auth, posts, comments, and reactions.
- Upgrades a connection to a WebSocket per authenticated user
  (`handlers/websocket.go`) for:
  - Broadcasting online/offline presence
  - Delivering new private messages in real time
  - Serving paginated chat history on request
  - Forcing logout on session invalidation
- Uses Go routines and channels to manage concurrent WebSocket clients.

### Frontend

- `routeer.js` implements client-side routing so the whole app lives on one
  HTML page.
- `ChatData.js` owns WebSocket connection lifecycle, chat/user state, and
  the socket message handlers.
- `Chatui.js` owns all chat-related DOM rendering (user list, message panel,
  message bubbles).
- `app-events.js` centralizes all DOM event delegation (`data-action`
  attributes) so behavior isn't scattered across inline handlers.

## Database

SQLite stores users, sessions, posts, categories, comments, reactions, and
private messages. See `backend/database/schema.sql` (and `schema.erd`) for
the full schema.

## Allowed Packages

- All standard Go packages
- [`gorilla/websocket`](https://github.com/gorilla/websocket)
- `sqlite3` (e.g. `mattn/go-sqlite3`)
- `golang.org/x/crypto/bcrypt`
- A UUID package (e.g. `google/uuid`)

No frontend libraries or frameworks are used.

## Getting Started

### Prerequisites

- Go (version matching `go.mod`)
- A C compiler (required by `mattn/go-sqlite3`, if used)

### Run

```bash
go run main.go
```

Then open the app in your browser at:

```
http://localhost:8082
```

(Adjust the port to match your `routing.go` / server configuration.)

## Learning Objectives

This project covers:

- HTML, HTTP, sessions & cookies, CSS
- Backend/Frontend separation and the DOM
- Go routines and channels for concurrency
- WebSockets, both server-side (Go) and client-side (JS)
- SQL and database manipulation
- Building a single-page application without a frontend framework

---

## Authors

- halhyane - [Gitea](https://learn.zone01oujda.ma/git/halhyane)
- elmehdi rezoug - [Github](https://github.com/elmehdi-rezoug)
