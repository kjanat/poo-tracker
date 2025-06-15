import { render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { describe, it, expect, vi } from "vitest";
import { ProtectedRoute } from "./ProtectedRoute";

const useAuthStoreMock = vi.fn();
vi.mock("../stores/authStore", () => ({
  useAuthStore: () => useAuthStoreMock()
}));

describe("ProtectedRoute", () => {
  it("redirects to login when not authenticated", () => {
    useAuthStoreMock.mockReturnValue({ isAuthenticated: false });
    render(
      <MemoryRouter initialEntries={["/secret"]}>
        <Routes>
          <Route
            path="/secret"
            element={
              <ProtectedRoute>
                <div>secret</div>
              </ProtectedRoute>
            }
          />
          <Route path="/login" element={<div>login page</div>} />
        </Routes>
      </MemoryRouter>
    );
    expect(screen.getByText(/login page/i)).toBeInTheDocument();
  });

  it("renders children when authenticated", () => {
    useAuthStoreMock.mockReturnValue({ isAuthenticated: true });
    render(
      <MemoryRouter initialEntries={["/secret"]}>
        <Routes>
          <Route
            path="/secret"
            element={
              <ProtectedRoute>
                <div>secret</div>
              </ProtectedRoute>
            }
          />
          <Route path="/login" element={<div>login page</div>} />
        </Routes>
      </MemoryRouter>
    );
    expect(screen.getByText(/secret/i)).toBeInTheDocument();
  });
});
