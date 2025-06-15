import { render } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import Logo from "./Logo";

describe("Logo component", () => {
  it("renders with default size", () => {
    const { container } = render(<Logo />);
    const svg = container.querySelector("svg");
    expect(svg?.getAttribute("width")).toBe("24");
    expect(svg?.getAttribute("height")).toBe("24");
  });

  it("applies custom size and class", () => {
    const { container } = render(<Logo size={48} className="custom" />);
    const svg = container.querySelector("svg") as SVGElement;
    expect(svg.getAttribute("width")).toBe("48");
    expect(svg.classList.contains("custom")).toBe(true);
  });
});
