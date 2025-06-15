import { afterEach, describe, expect, it, vi } from "vitest";
import { useAuthStore } from "./authStore";

const userResponse = {
  token: "tok",
  user: { id: "1", email: "a@b.com", name: "Tester" }
};

describe("authStore", () => {
  afterEach(() => {
    useAuthStore.getState().logout();
    vi.restoreAllMocks();
  });

  it("logs in and updates state", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(() =>
        Promise.resolve({
          ok: true,
          json: () => Promise.resolve(userResponse)
        })
      ) as any
    );

    await useAuthStore.getState().login("a@b.com", "secret");
    const state = useAuthStore.getState();
    expect(state.isAuthenticated).toBe(true);
    expect(state.user?.email).toBe("a@b.com");
    expect(state.token).toBe("tok");
  });
});
