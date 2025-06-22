# 💩 Copilot Vibes: Poo Tracker

## Repo: [github.com/kjanat/poo-tracker](https://github.com/kjanat/poo-tracker)

### TL;DR

Keep it clean, make it work, and try not to break the toilet.

---

## 🛠️ What We’re Using

- **Frontend**: React (Vite + TypeScript + TailwindCSS)
- **Backend**: Go with Gin – clean arch, in-memory for now, JWT-auth’d
- **DB**: PostgreSQL planned, in-memory for dev
- **File Storage**: S3-compatible (MinIO locally)
- **AI Sidekick**: Python FastAPI in Docker, accessed via Go backend
- **CI/CD**: GitHub Actions
- **Testing**: Vitest, go test, pytest
- **Infra**: Everything Dockerized
- **Package Managers**: `pnpm` (TS) and `uv` (Python)

---

## 🧘 Chill Dev Guidelines

- **TypeScript all the way** – plain JS just makes stuff harder later
- **Keep components bite-sized** – if something’s getting chunky, break it up
- **Tailwind only** – no global CSS drama, no modules, just utility classes
- **REST over everything** – no GraphQL unless you’ve got a damn good reason
- **Go patterns** – handle your errors, use pointers when it makes sense, follow the repo vibe
- **Streaming uploads** – no base64 blobs in the DB, we’re not monsters
- **Python lives behind Go** – don’t let the frontend chat with the AI directly
- **Write tests** – doesn’t have to be perfect, just don’t leave future-you crying
- **Use `.env.example`** – no real keys in the repo, ever
- **Pre-commit hooks** – formatters and linters keep the peace

---

## 🔀 Branching & PRs

- Branch off `master`
- Name it like: `feat/shiny-feature`, `fix/bug`, etc.
- PRs should:
  - Have a short description
  - Pass tests
  - Get at least one sleepy eyeball to review
  - Not break CI

---

## 🙃 Stuff to Avoid

- No direct commits to `master`
- No random packages unless you know what they’re doing
- No leftover `console.log` chaos
- No TODOs without context

---

## 🧻 Final Words

- If it looks like 💩, smells like 💩, or compiles like 💩... just refactor it.
- Clean code = clean colon.
- Don’t be a jerk. Review like you’d want to be reviewed.

---

Now go commit something beautiful.
