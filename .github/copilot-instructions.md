# Copilot Instructions: Poo Tracker

## Repo: [github.com/kjanat/poo-tracker](https://github.com/kjanat/poo-tracker)

### TL;DR

If your code stinks, it won’t get merged. And yes, we’ll know.

---

## Tech Stack

- **Frontend**: React (with Vite), TypeScript, TailwindCSS
- **Backend/API**: Go (Gin framework, dependency injection, in-memory repositories)
  - Clean architecture with repository pattern
  - JWT authentication with user management
  - Enhanced BowelMovement model with separated details for performance
  - Comprehensive validation and error handling
- **Database**: PostgreSQL (planned - currently using in-memory for development)
- **Storage**: S3-compatible for turd-photos (MinIO for local, AWS S3 for prod)
- **AI / ML**: Python FastAPI service (Dockerized), talks to Go backend via REST
- **CI/CD**: GitHub Actions
- **Testing**: Vitest (frontend), go test (backend), pytest (AI service)
- **Linting**: ESLint with @typescript-eslint and Prettier (don't fight the config), golangci-lint (backend)
- **Infra**: Docker everywhere. No raw installs.
- **Package Managers**:
  - pnpm (because npm is for amateurs)
  - uv (for Python, because pip is for amateurs)

---

## Coding Practices

1. **TypeScript is not optional.**

   If you sneak in any plain JavaScript, I will revert your commit and send you a brown paper bag.

2. **Component-based everything.**

   No 1,000-line files. If your React component grows a tumor, split it.

3. **Tailwind > CSS files**

   No CSS modules, no SCSS, no “quick fixes” in styles.css. Utility-first or GTFO.

4. **Backend: RESTful.**

   If you say "GraphQL" I'll call your mother.

5. **Go backend standards:**

   - Always handle error return values (errcheck will catch you)
   - Use pointer types for optional fields to distinguish nil from zero values
   - Repository pattern with clean interfaces
   - Comprehensive validation for all inputs
   - JWT middleware for protected endpoints

6. **Database migrations happen automatically.**

   No SQL scripts in the repo, unless you're writing a damn seed file.

7. **Pre-commit hooks for code quality.**

   All linting, formatting, and type-checking is managed via [pre-commit](https://pre-commit.com) and the `.pre-commit-config.yaml` in the project root. Husky and lint-staged are no longer used.

   - Install with `uv tool install pre-commit` (or `pip install pre-commit`)
   - Run `pre-commit install` after cloning
   - Hooks run on every commit, or manually with `pre-commit run --all-files`

8. **Photo uploads must stream.**

   No base64 blobs in the database, you animal.

9. **AI service is firewalled from main app.**

   Only ever accessible via backend. Don't let your SPA talk directly to Python, or you'll make me regret giving you network access.

10. **Write tests.**

    Every new feature needs tests. If you "forget", the next person gets to refactor your code, and nobody wants that.

11. **Env vars go in `.env.example`**

Don't hardcode keys, don't commit real secrets.

12. **Follow commit conventions**:

    ```txt
    feat: add Bristol stool chart selector
    fix: correct off-by-one in shit counter
    chore: update dependencies
    ```

---

## Branching, PRs, and CI

- **Branch off `master`**
- Use feature branches (`feat/`, `fix/`, `chore/`).
- Every PR needs a description, a passing test suite, and a code review from someone who’s at least 80% awake.
- Green CI/CD or no merge.

---

## What NOT to do

- No direct pushes to `master`
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
