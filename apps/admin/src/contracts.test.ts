import { describe, expect, it } from "vitest";
import { FOUNDATION_AREAS, PRESENTATION_LAYERS } from "./contracts";

describe("WP-001 admin foundation contracts", () => {
  it("keeps the four accepted presentation layers explicit", () => {
    expect(PRESENTATION_LAYERS).toHaveLength(4);
    expect(new Set(PRESENTATION_LAYERS).size).toBe(4);
  });

  it("keeps each foundation area owned by a repository boundary", () => {
    expect(FOUNDATION_AREAS.every((area) => area.owner.length > 0)).toBe(true);
    expect(FOUNDATION_AREAS.every((area) => area.state === "foundation")).toBe(true);
  });
});
