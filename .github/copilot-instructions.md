Here you go, Kaj. **copilot-instructions.md**—no-nonsense, clear, direct, and fully tailored for the average dev who’s about to cry because they have to touch something called “Poo Tracker”. No corporate cringe, just the brutal reality of what you want (and what you absolutely don’t).

---

# Copilot Instructions: Poo Tracker

## Repo: [github.com/kjanat/poo-tracker](https://github.com/kjanat/poo-tracker)

### TL;DR

If your code stinks, it won’t get merged. And yes, we’ll know.

---

## Tech Stack

- **Frontend**: React (with Vite), TypeScript, TailwindCSS
- **Backend/API**: Node.js, Express, TypeScript
- **Database**: PostgreSQL (use Prisma ORM, or GTFO)
- **Storage**: S3-compatible for turd-photos (MinIO for local, AWS S3 for prod)
- **AI / ML**: Python FastAPI service (Dockerized), talks to Node backend via REST
- **CI/CD**: GitHub Actions
- **Testing**: Vitest (frontend), Jest/Supertest (backend), pytest (AI service)
- **Linting**: ESLint (StandardJS rules), Prettier (and don’t “fix” the config, just follow it)
- **Infra**: Docker everywhere. No raw installs.

---

## Coding Practices

1. **TypeScript is not optional.**
   If you sneak in any plain JavaScript, I will revert your commit and send you a brown paper bag.
2. **Component-based everything.**
   No 1,000-line files. If your React component grows a tumor, split it.
3. **Tailwind > CSS files**
   No CSS modules, no SCSS, no “quick fixes” in styles.css. Utility-first or GTFO.
4. **Backend: RESTful.**
   If you say “GraphQL” I’ll call your mother.
5. **Database migrations with Prisma only.**
   No SQL scripts in the repo, unless you’re writing a damn seed file.
6. **Photo uploads must stream.**
   No base64 blobs in the database, you animal.
7. **AI service is firewalled from main app.**
   Only ever accessible via backend. Don’t let your SPA talk directly to Python, or you’ll make me regret giving you network access.
8. **Write tests.**
   Every new feature needs tests. If you “forget”, the next person gets to refactor your code, and nobody wants that.
9. **Env vars go in `.env.example`**
   Don’t hardcode keys, don’t commit real secrets.
10. **Follow commit conventions**:

    ```
    feat: add Bristol stool chart selector
    fix: correct off-by-one in shit counter
    chore: update dependencies
    ```

---

## Branching, PRs, and CI

- **Branch off `main`**
- Use feature branches (`feat/`, `fix/`, `chore/`).
- Every PR needs a description, a passing test suite, and a code review from someone who’s at least 80% awake.
- Green CI/CD or no merge.

---

## What NOT to do

- No direct pushes to `main`
- No wild dependencies (“left-pad” jokes will get you a one-way ticket to dependency hell)
- No copying Stack Overflow code without understanding it
- No “console.log debugging” left in PRs
- No “fix later” TODOs, unless it’s properly documented and assigned

---

## Golden Rules

- If you wouldn’t eat it, don’t ship it.
- Clean code = Clean colon.
- Don’t be an asshole in code reviews, unless the code really deserves it.

---

That’s it. Now go code, you magnificent bastard.
