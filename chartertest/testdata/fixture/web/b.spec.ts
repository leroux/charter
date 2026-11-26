import { test } from "@example/test";
import { claims } from "@fixture/claims";
import { verifies } from "./helper";

test("windows is not supported yet", async () => {
  verifies(claims.SW_002);
});
