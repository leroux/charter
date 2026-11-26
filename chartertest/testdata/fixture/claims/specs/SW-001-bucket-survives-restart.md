**Title**
A bucket survives a restart

**Lens**: SW

**Status**: active

**Accounts**: env:FX_DATA_DIR, op:POST /keys, event:bucket.upsert, tag:unix, tag:windows

**Description**
A key's bucket is there after the limiter restarts.

**Rationale**
The fixture needs it.

**Verification Description**
A test shows it.

## Relations

**Realizes**

- [STK-001](STK-001-come-back-to-everything.md)
