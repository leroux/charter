# Needs

What a caller of a service behind ratelimit, and the operator who runs it,
can rely on. A need is verified from their side: the ratelimit command,
fed keys and read back, as a script in front of a service would drive it.

The key words MUST, MUST NOT, SHOULD, SHOULD NOT and MAY in these
documents are to be interpreted as described in BCP 14 (RFC 2119, RFC 8174)
when, and only when, they appear in all capitals.

## STK-261126-H2SW ratelimit MUST let a quiet caller burst, then hold it to its sustained rate, and MUST NOT let one caller spend another's requests

Status: active

- ratelimit MUST let a caller that has been quiet make a burst of
  requests at once.
- After the burst, ratelimit MUST let the caller through at its sustained
  rate and MUST deny a request beyond either.
- One caller's requests MUST NOT use up another caller's.

Why: a service behind ratelimit needs protection from any one caller that
sends more than its share, and a caller that keeps to its limits needs to
keep getting through while a noisy one sends too much. Limiting each
caller on its own meets both: the noisy caller is denied, and every other
caller keeps its burst and its rate.

Scenarios:
1. `burst-then-deny` (agent, provisional): given the default burst of
   five, when a quiet caller sends six requests at once, then the first
   five are allowed and the sixth is denied.
2. `callers-apart` (agent, provisional): given a caller that has spent
   its burst and been denied, when another caller sends a request at the
   same instant, then that request is allowed.

## STK-261126-T2BD ratelimit MUST take its rate and burst from the operator at start, and MUST NOT start on a setting it cannot read

Status: active

- ratelimit MUST let the operator set the sustained rate and the burst
  when starting it, without rebuilding it.
- A setting ratelimit cannot read MUST stop it at the start, with a
  failing exit and the setting named, so that it never limits by a value
  the operator did not choose.

Why: the right limits depend on the deployment, so the operator sets them
when starting ratelimit rather than when building it. A limiter that fell
back to a default when a setting was mistyped would hide the typo until
the service behind it fell over. Stopping at the start shows the typo
while the operator is there to fix it, at the cost of a start that fails
where it could have run on a guess.

Scenarios:
1. `burst-set` (agent, provisional): given the burst set to two, when a
   caller sends three requests at once, then the third is denied.
2. `bad-setting-stops` (agent, provisional): given the rate set to
   `fast`, when ratelimit starts with a key on its input, then it exits
   with a failing status, gives no verdict and names `RATELIMIT_RATE`.
