import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, it, expect, vi } from "vitest";
import { HomePage } from "./HomePage";

const useAuthStoreMock = vi.fn();
vi.mock("../stores/authStore", () => ({
  useAuthStore: () => useAuthStoreMock()
}));

describe("HomePage", () => {
  it("shows call to action when not authenticated", () => {
    useAuthStoreMock.mockReturnValue({ isAuthenticated: false });
    render(
      <MemoryRouter>
        <HomePage />
      </MemoryRouter>
    );
    expect(screen.getByText(/Get Started/i)).toBeInTheDocument();
  });

  it("shows dashboard link when authenticated", () => {
    useAuthStoreMock.mockReturnValue({ isAuthenticated: true });
    render(
      <MemoryRouter>
        <HomePage />
      </MemoryRouter>
    );
    expect(screen.getByText(/Go to Dashboard/i)).toBeInTheDocument();
  });
});
