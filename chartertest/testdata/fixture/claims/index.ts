// Generated from claims by the charter gate. Do not edit.
// Regenerate: CHARTER_UPDATE_TYPESCRIPT=1 go test -count=1 ./claims

export const claims = {
  // STK
  /** A person can come back and find everything as it was */
  STK_001: "STK-001",
  /** Planned: An operator can run the limiter on another machine */
  STK_002: "STK-002",

  // SW
  /** A bucket survives a restart */
  SW_001: "SW-001",
  /** @deprecated The limiter runs on Windows */
  SW_002: "SW-002",
  /** Planned: A person can clear a quota's counters */
  SW_003: "SW-003",

  // CON
  /** The lock file is removed on close */
  CON_001: "CON-001",
  /** The fold never panics */
  CON_002: "CON-002",
} as const;

/** A spec id: what a test verifies. */
export type Claim = (typeof claims)[keyof typeof claims];
