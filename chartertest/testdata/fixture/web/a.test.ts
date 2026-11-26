import { describe, it } from "@example/test";
import { claims as c } from "@fixture/claims";
import { verifies } from "./helper";

describe("bucket", () => {
  it("survives a restart", () => {
    verifies(c.SW_001);
  });
  it.each([1, 2])("folds %d", () => {
    verifies(c.CON_002);
  });
  it(`locks ${"maybe"}`, () => {
    verifies(c.CON_001);
  });
});
