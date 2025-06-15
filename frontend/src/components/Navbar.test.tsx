import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { useAuthStore } from "../stores/authStore";
import { Navbar } from "./Navbar";
import { vi, describe, it, expect } from "vitest";

describe("Navbar", () => {
  it("shows login link when unauthenticated", () => {
    useAuthStore.setState({ user: null, token: null, isAuthenticated: false });
    render(
      <MemoryRouter>
        <Navbar />
      </MemoryRouter>
    );
    expect(screen.getByText(/login/i)).toBeInTheDocument();
  });

  it("shows dashboard link when authenticated", () => {
    useAuthStore.setState({
      user: { id: "1", email: "a@b.c", name: "A" },
      token: "tok",
      isAuthenticated: true,
      login: useAuthStore.getState().login,
      setAuth: useAuthStore.getState().setAuth,
      logout: vi.fn()
    });
    render(
      <MemoryRouter>
        <Navbar />
      </MemoryRouter>
    );
    expect(screen.getAllByText(/dashboard/i)[0]).toBeInTheDocument();
  });
});
