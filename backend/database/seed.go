package database

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// DemoPassword is the shared password for every seeded demo account.
// The quick-login buttons on the login page use it too — keep it in sync
// with frontend/static/js/pages/login.js (DEMO_PASSWORD).
const DemoPassword = "demo1234"

type seedUser struct {
	Nickname  string
	FirstName string
	LastName  string
	Age       int
	Gender    string
	Email     string
	// HoursAgo sets how long ago the account was last seen (sidebar ordering).
	HoursAgo float64
}

type seedPost struct {
	Author     int // index into seedUsers
	Title      string
	Content    string
	Categories []int // category IDs (1 General, 2 Programming, 3 Gaming, 4 Movies, 5 Sports, 6 Anime)
}

type seedComment struct {
	Post     int // index into seedPosts
	Author   int // index into seedUsers
	Content  string
	HoursAgo float64
}

type seedReaction struct {
	Target int // post or comment index
	User   int // index into seedUsers
	Like   bool
}

type seedMessage struct {
	From     int // index into seedUsers
	To       int // index into seedUsers
	Content  string
	HoursAgo float64
}

var seedUsers = []seedUser{
	{Nickname: "Sara", FirstName: "Sara", LastName: "Johnson", Age: 26, Gender: "female", Email: "sara@example.com", HoursAgo: 4},
	{Nickname: "Ahmed", FirstName: "Ahmed", LastName: "Ali", Age: 24, Gender: "male", Email: "ahmed@example.com", HoursAgo: 20},
	{Nickname: "Nicole", FirstName: "Nicole", LastName: "Martin", Age: 29, Gender: "female", Email: "nicole@example.com", HoursAgo: 16},
	{Nickname: "Alex", FirstName: "Alex", LastName: "Brown", Age: 31, Gender: "male", Email: "alex@example.com", HoursAgo: 5},
}

var seedPosts = []seedPost{
	{
		Author: 0, Title: "Welcome to 01Forum 🎉", Categories: []int{1},
		Content: "Hey everyone, welcome to 01Forum!\n\n" +
			"This little community is a fully real-time single-page app: posts, comments, reactions, and private chat all update instantly over WebSockets — no page reloads.\n\n" +
			"A few things to try:\n" +
			"• Browse posts and filter them by category\n" +
			"• Open a post, drop a comment, like or dislike\n" +
			"• Click a user in the sidebar and send them a message — typing indicators included\n\n" +
			"Have fun, stay kind, and happy posting! 💬",
	},
	{
		Author: 0, Title: "Go vs Rust for backend APIs in 2026", Categories: []int{2},
		Content: "I've built small backends in both Go and Rust this year, and honestly both are great — it comes down to what you optimize for.\n\n" +
			"Go wins on developer velocity: instant compiles, simple concurrency with goroutines, and a stdlib that covers 90% of what you need for an API. Deploying a single static binary is effortless.\n\n" +
			"Rust wins when you need maximum performance and memory safety guarantees — at the cost of fighting the borrow checker. For a chat app like this one, Go's goroutine-per-connection model is hard to beat.\n\n" +
			"My rule of thumb: if it's an I/O-bound web service, Go. If it's a CPU-bound system component, Rust. What do you all think?",
	},
	{
		Author: 1, Title: "SQLite is criminally underrated", Categories: []int{2},
		Content: "Everyone reaches for Postgres out of habit, but SQLite deserves way more love for side projects and small apps.\n\n" +
			"It's a full relational database in one file, zero configuration, and it's blindingly fast for reads. This very forum runs on SQLite and it handles everything fine.\n\n" +
			"Modern SQLite even has WAL mode, so concurrent readers and writers don't block each other.\n\n" +
			"If your app fits on one machine, don't over-engineer it. Start with SQLite.",
	},
	{
		Author: 1, Title: "Best open-world games to sink 100 hours into", Categories: []int{3},
		Content: "I've been rotating through open worlds lately and these are the ones that actually respect your time:\n\n" +
			"• Elden Ring — brutal, but every discovery feels earned\n" +
			"• Zelda: Tears of the Kingdom — pure creative freedom\n" +
			"• Red Dead Redemption 2 — the most alive world ever made\n" +
			"• Cyberpunk 2077 — fully redeemed after the 2.0 update\n\n" +
			"What's your favorite open world? I need something new to disappear into.",
	},
	{
		Author: 1, Title: "Hot take: indie games beat AAA releases", Categories: []int{3},
		Content: "I'll say it: most AAA games this year were safe, bloated, and buggy at launch, while indie games are where the innovation actually lives.\n\n" +
			"Hollow Knight, Hades, Stardew Valley, Outer Wilds — some of the best games ever made, all built by tiny teams with big ideas.\n\n" +
			"AAA has the budget, but indie has the soul. Change my mind.",
	},
	{
		Author: 2, Title: "What are you watching this weekend?", Categories: []int{4},
		Content: "Weekend movie night planning! I'm thinking of doing a double feature.\n\n" +
			"On the list: a slow-burn thriller, a comfort comedy, and maybe a rewatch of an old classic.\n\n" +
			"What's on your screen this weekend? Give me your best recommendations — anything except horror, I'm a scaredy-cat 🫣",
	},
	{
		Author: 2, Title: "Underrated sci-fi movies you need to see", Categories: []int{4},
		Content: "Everyone knows the classics, so here are my underrated sci-fi picks that deserve more attention:\n\n" +
			"• Arrival — the most human alien movie ever made\n" +
			"• Ex Machina — a tight, tense chamber piece\n" +
			"• Coherence — mind-bending, made for almost nothing\n" +
			"• Annihilation — beautiful and deeply unsettling\n\n" +
			"I've watched each of these at least three times. What's the most underrated sci-fi you've seen?",
	},
	{
		Author: 3, Title: "Champions League final predictions 🏆", Categories: []int{5},
		Content: "The final is this weekend and I genuinely can't call it.\n\n" +
			"Both teams are in incredible form and the head-to-head this season is dead even. The midfield battle will decide it, in my opinion — whoever controls the tempo in the first half wins.\n\n" +
			"My prediction: 2-1, a late winner, and my heart rate not recovering for a week. Drop your predictions below!",
	},
	{
		Author: 3, Title: "Marathon training for total beginners", Categories: []int{5},
		Content: "I get asked about marathon training a lot, so here's the beginner plan that worked for me:\n\n" +
			"1. Follow a plan — Couch to 5K, then a half, then a full. Slow progress beats heroics.\n" +
			"2. Long runs on weekends — one long run a week, adding about 10% distance each time.\n" +
			"3. Rest days are training days. Your body adapts while recovering.\n" +
			"4. Nutrition and sleep matter more than any gear. Get decent shoes and stop overthinking.\n\n" +
			"It took me two years from zero to the finish line. You can do it!",
	},
	{
		Author: 2, Title: "Anime with the best soundtracks", Categories: []int{6},
		Content: "A great soundtrack elevates a good show into an unforgettable one. My all-time favorites:\n\n" +
			"• Cowboy Bebop — jazz that defines the entire vibe\n" +
			"• Your Name — RADWIMPS made every scene ache\n" +
			"• Attack on Titan — epic orchestral bangers\n" +
			"• Violet Evergarden — pure emotional devastation\n\n" +
			"Drop your favorite anime OSTs below, I'm building a playlist.",
	},
	{
		Author: 0, Title: "Why WebSockets beat polling for chat apps", Categories: []int{2, 1},
		Content: "If you're building anything real-time, do yourself a favor and skip polling.\n\n" +
			"Polling wastes bandwidth, adds latency, and makes your server do work even when nothing happened. WebSockets keep one persistent connection open and push messages the instant they're available.\n\n" +
			"This forum's chat uses a WebSocket hub with a goroutine per connection and a mutex-protected registry — plus ping/pong keepalives so dead connections get cleaned up.\n\n" +
			"Latency went from 'refresh to check' to 'instant'. The difference feels magical.",
	},
	{
		Author: 3, Title: "E-sports are real sports. Change my mind.", Categories: []int{3, 5},
		Content: "I keep hearing 'e-sports aren't real sports' and I think it's time we talked.\n\n" +
			"E-sports athletes train 8+ hours a day, follow strict routines, compete in front of millions, and win real championships. The physical demands — reaction time, hand-eye coordination, endurance — are measurable.\n\n" +
			"If chess is a sport, competitive gaming is a sport. You don't need a ball to be athletic. Change my mind.",
	},
}

var seedComments = []seedComment{
	{Post: 0, Author: 1, HoursAgo: 27, Content: "First! 🎉 Great intro, love the vibe of this place already."},
	{Post: 0, Author: 2, HoursAgo: 26, Content: "Happy to be here! The UI is gorgeous by the way 😍"},
	{Post: 0, Author: 3, HoursAgo: 25, Content: "Awesome community. The real-time chat is super smooth."},
	{Post: 1, Author: 1, HoursAgo: 24, Content: "Team Go all the way. Goroutines make my life so much easier."},
	{Post: 1, Author: 3, HoursAgo: 23, Content: "I'm learning Rust right now and... ouch. But it's worth it."},
	{Post: 2, Author: 0, HoursAgo: 22, Content: "Preach. WAL mode + foreign keys on and you're 90% of the way to Postgres for small apps."},
	{Post: 3, Author: 3, HoursAgo: 21, Content: "Elden Ring ruined other games for me. Nothing compares now."},
	{Post: 3, Author: 2, HoursAgo: 20, Content: "RDR2 is the only game that made me cry. That world is unreal."},
	{Post: 4, Author: 2, HoursAgo: 19, Content: "Hades 2 can't come soon enough. Agree 100%."},
	{Post: 5, Author: 0, HoursAgo: 18, Content: "Watching a cozy mystery tonight. Comfort over thrillers for me 🍿"},
	{Post: 5, Author: 3, HoursAgo: 17, Content: "Rewatching The Godfather for the fifth time. Never gets old."},
	{Post: 6, Author: 0, HoursAgo: 16, Content: "Arrival destroyed me. I still think about that ending weekly."},
	{Post: 6, Author: 1, HoursAgo: 15, Content: "Coherence is a masterpiece. So much bang for so little budget."},
	{Post: 7, Author: 1, HoursAgo: 14, Content: "2-1 is brave. I'm saying 1-1 and a penalty shootout — my nerves can't take it."},
	{Post: 8, Author: 2, HoursAgo: 13, Content: "Bookmarking this. 'Rest days are training days' — stealing that."},
	{Post: 9, Author: 1, HoursAgo: 12, Content: "Cowboy Bebop's 'Tank!' is the best opening ever made. I'll fight anyone on this."},
	{Post: 10, Author: 3, HoursAgo: 11, Content: "The instant chat on this site is seriously impressive. No refresh, no lag."},
	{Post: 11, Author: 1, HoursAgo: 10, Content: "If strategy counts, then chess and e-sports are both sports. Facts."},
}

// seedPostReactions: reactions on posts (index into seedPosts, user index, like bool).
var seedPostReactions = []seedReaction{
	{Target: 0, User: 1, Like: true},
	{Target: 0, User: 2, Like: true},
	{Target: 0, User: 3, Like: true},
	{Target: 1, User: 1, Like: true},
	{Target: 1, User: 3, Like: true},
	{Target: 2, User: 0, Like: true},
	{Target: 2, User: 3, Like: true},
	{Target: 3, User: 0, Like: true},
	{Target: 3, User: 2, Like: true},
	{Target: 3, User: 3, Like: true},
	{Target: 4, User: 0, Like: true},
	{Target: 4, User: 2, Like: true},
	{Target: 5, User: 0, Like: true},
	{Target: 5, User: 3, Like: true},
	{Target: 6, User: 0, Like: true},
	{Target: 6, User: 1, Like: true},
	{Target: 7, User: 1, Like: true},
	{Target: 8, User: 0, Like: true},
	{Target: 8, User: 2, Like: true},
	{Target: 9, User: 0, Like: true},
	{Target: 9, User: 1, Like: true},
	{Target: 10, User: 1, Like: true},
	{Target: 10, User: 2, Like: true},
	{Target: 10, User: 3, Like: true},
	{Target: 11, User: 1, Like: true},
	{Target: 11, User: 2, Like: false},
}

// seedCommentReactions: reactions on comments (index into seedComments, user index, like bool).
var seedCommentReactions = []seedReaction{
	{Target: 0, User: 2, Like: true},
	{Target: 0, User: 3, Like: true},
	{Target: 1, User: 0, Like: true},
	{Target: 1, User: 3, Like: true},
	{Target: 6, User: 1, Like: true},
	{Target: 6, User: 2, Like: true},
	{Target: 9, User: 3, Like: true},
	{Target: 11, User: 1, Like: true},
	{Target: 11, User: 2, Like: true},
	{Target: 14, User: 0, Like: true},
	{Target: 14, User: 3, Like: true},
	{Target: 15, User: 0, Like: true},
	{Target: 15, User: 2, Like: true},
}

// seedMessages: private conversations so the chat sidebar looks alive on first login.
var seedMessages = []seedMessage{
	// Sara ↔ Ahmed — SQLite chat
	{From: 0, To: 1, HoursAgo: 30, Content: "Hey! Saw you mention SQLite — here's the trick for chat apps: enable WAL mode."},
	{From: 1, To: 0, HoursAgo: 29.8, Content: "WAL mode? Never heard of it. What does that actually change?"},
	{From: 0, To: 1, HoursAgo: 29.5, Content: "Writers don't block readers anymore. Huge for a real-time app like this one 😄"},
	{From: 1, To: 0, HoursAgo: 29, Content: "Oh nice, going to try that tonight. Thanks a lot!"},
	{From: 0, To: 1, HoursAgo: 28.5, Content: "Anytime! Ping me if you hit any issues."},
	// Nicole ↔ Sara — movies
	{From: 2, To: 0, HoursAgo: 26, Content: "Have you seen the new sci-fi movie that just dropped on streaming? It's incredible."},
	{From: 0, To: 2, HoursAgo: 25.5, Content: "Not yet! Is it the one everyone's talking about? Adding it to my list."},
	{From: 2, To: 0, HoursAgo: 25, Content: "Yes! You'll love it — right up your alley. Let me know what you think."},
	{From: 0, To: 2, HoursAgo: 24, Content: "Just finished it. WOW. That ending... 😱"},
	// Alex ↔ Ahmed — gaming
	{From: 3, To: 1, HoursAgo: 22, Content: "That indie game you recommended is amazing. 10/10, would recommend right back."},
	{From: 1, To: 3, HoursAgo: 21.5, Content: "Told you! Indie devs are absolutely killing it right now."},
	{From: 3, To: 1, HoursAgo: 21, Content: "The soundtrack alone is worth it. Anyway, up for a match this weekend?"},
	{From: 1, To: 3, HoursAgo: 20, Content: "Definitely. Saturday afternoon?"},
	// Nicole ↔ Alex — running
	{From: 2, To: 3, HoursAgo: 18, Content: "Booked my first half marathon! Any training advice?"},
	{From: 3, To: 2, HoursAgo: 17.5, Content: "Congrats! Start slow, don't skip rest days, and get good shoes."},
	{From: 2, To: 3, HoursAgo: 17, Content: "Thanks! I'll take that advice. Rest days it is 😅"},
	{From: 3, To: 2, HoursAgo: 16.5, Content: "You've got this. See you at the finish line!"},
	// Alex ↔ Sara — welcome
	{From: 3, To: 0, HoursAgo: 5, Content: "The real-time chat in this forum is so smooth. Great work!"},
	{From: 0, To: 3, HoursAgo: 4.5, Content: "Thanks! All Go + WebSockets under the hood, no magic 🪄"},
}

// Seed inserts demo data on first startup (only when the users table is empty),
// so it never duplicates or clobbers real accounts. It is safe to call on
// every boot: once data exists it simply logs and returns.
func Seed() error {
	tx, err := Database.Begin()
	if err != nil {
		return fmt.Errorf("seed: begin tx: %w", err)
	}
	defer tx.Rollback() // no-op after Commit

	// Idempotency guard runs inside the transaction so two concurrent
	// startups can't both pass it and double-seed.
	var count int
	if err := tx.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return fmt.Errorf("seed: count users: %w", err)
	}
	if count > 0 {
		log.Println("Seed: demo data already present, skipping")
		return nil
	}

	now := time.Now()

	// 1. Users (bcrypt-hashed, shared DemoPassword).
	userIDs := make([]int64, len(seedUsers))
	for i, u := range seedUsers {
		hash, err := bcrypt.GenerateFromPassword([]byte(DemoPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("seed: hash password for %s: %w", u.Nickname, err)
		}
		lastSeen := now.Add(-time.Duration(u.HoursAgo * float64(time.Hour))).Format("2006-01-02 15:04:05")
		res, err := tx.Exec(
			`INSERT INTO users (nickname, first_name, last_name, age, gender, email, password, last_seen)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			u.Nickname, u.FirstName, u.LastName, u.Age, u.Gender, u.Email, string(hash), lastSeen,
		)
		if err != nil {
			return fmt.Errorf("seed: insert user %s: %w", u.Nickname, err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("seed: user id for %s: %w", u.Nickname, err)
		}
		userIDs[i] = id
	}

	// 2. Posts + their categories.
	postIDs := make([]int64, len(seedPosts))
	for i, p := range seedPosts {
		res, err := tx.Exec(
			"INSERT INTO posts (title, content, user_id) VALUES (?, ?, ?)",
			p.Title, p.Content, userIDs[p.Author],
		)
		if err != nil {
			return fmt.Errorf("seed: insert post %d: %w", i, err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("seed: post id %d: %w", i, err)
		}
		postIDs[i] = id
		for _, cat := range p.Categories {
			if _, err := tx.Exec(
				"INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)",
				id, cat,
			); err != nil {
				return fmt.Errorf("seed: post %d category %d: %w", i, cat, err)
			}
		}
	}

	// 3. Comments (staggered timestamps).
	commentIDs := make([]int64, len(seedComments))
	for i, c := range seedComments {
		createdAt := now.Add(-time.Duration(c.HoursAgo * float64(time.Hour))).Format("2006-01-02 15:04:05")
		res, err := tx.Exec(
			"INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
			postIDs[c.Post], userIDs[c.Author], c.Content, createdAt,
		)
		if err != nil {
			return fmt.Errorf("seed: insert comment %d: %w", i, err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("seed: comment id %d: %w", i, err)
		}
		commentIDs[i] = id
	}

	// 4. Reactions on posts.
	for _, r := range seedPostReactions {
		isLike := 0
		if r.Like {
			isLike = 1
		}
		if _, err := tx.Exec(
			"INSERT INTO post_reactions (user_id, post_id, is_like) VALUES (?, ?, ?)",
			userIDs[r.User], postIDs[r.Target], isLike,
		); err != nil {
			return fmt.Errorf("seed: post reaction: %w", err)
		}
	}

	// 5. Reactions on comments.
	for _, r := range seedCommentReactions {
		isLike := 0
		if r.Like {
			isLike = 1
		}
		if _, err := tx.Exec(
			"INSERT INTO comment_reactions (user_id, comment_id, is_like) VALUES (?, ?, ?)",
			userIDs[r.User], commentIDs[r.Target], isLike,
		); err != nil {
			return fmt.Errorf("seed: comment reaction: %w", err)
		}
	}

	// 6. Private messages (staggered timestamps, ms precision like the column default).
	for i, m := range seedMessages {
		createdAt := now.Add(-time.Duration(m.HoursAgo * float64(time.Hour))).Format("2006-01-02 15:04:05.000")
		if _, err := tx.Exec(
			"INSERT INTO messages (sender_id, receiver_id, content, created_at) VALUES (?, ?, ?, ?)",
			userIDs[m.From], userIDs[m.To], m.Content, createdAt,
		); err != nil {
			return fmt.Errorf("seed: insert message %d: %w", i, err)
		}
	}

	// 7. Derive last_seen from the seeded activity so the sidebar's presence
	// ordering stays consistent with the conversations above.
	if _, err := tx.Exec(
		`UPDATE users SET last_seen = (
			SELECT MAX(m.created_at) FROM messages m
			WHERE m.sender_id = users.id OR m.receiver_id = users.id
		)`,
	); err != nil {
		return fmt.Errorf("seed: derive last_seen: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("seed: commit: %w", err)
	}

	log.Printf("Seed: inserted %d users, %d posts, %d comments, %d messages",
		len(seedUsers), len(seedPosts), len(seedComments), len(seedMessages))
	return nil
}
