import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";
import { useAuthStore } from "./stores/authStore";

const originalFetch = global.fetch;

describe("useAuthStore", () => {
  beforeEach(() => {
    global.fetch = vi.fn();
    useAuthStore.setState({ user: null, token: null, isAuthenticated: false });
  });

  afterEach(() => {
    global.fetch = originalFetch;
  });

  it("login stores token and user", async () => {
    (global.fetch as any).mockResolvedValue({
      ok: true,
      json: () =>
        Promise.resolve({
          user: { id: "1", email: "test@example.com" },
          token: "tok"
        })
    });

    await useAuthStore.getState().login("test@example.com", "password");

    const state = useAuthStore.getState();
    expect(state.isAuthenticated).toBe(true);
    expect(state.token).toBe("tok");
    expect(state.user?.email).toBe("test@example.com");
  });

  it("login throws on failure", async () => {
    (global.fetch as any).mockResolvedValue({
      ok: false,
      json: () => Promise.resolve({ error: "fail" })
    });

    await expect(
      useAuthStore.getState().login("bad@example.com", "wrong")
    ).rejects.toThrow("fail");
  });
});
