# Needs

What a caller of a service behind ratelimit, and the operator who runs it,
can rely on. A need is verified from their side: the ratelimit command,
fed keys and read back, as a script in front of a service would drive it.

## STK-261126-H2SW A caller gets a burst at once, then its sustained rate, and no more

Status: active

A caller that has been quiet may make a burst of requests at once; after
that it is let through at its sustained rate, and a request beyond either
is denied. One caller's requests never use up another's.

Why: a service behind ratelimit is protected from any one caller, and a
well-behaved caller is never punished for a noisy one.

Scenarios:
1. `burst-then-deny` (agent, provisional): a quiet caller that sends one
   more request than the burst at once is let through for the burst and
   denied for the last.
2. `callers-apart` (agent, provisional): a caller that has spent its
   burst does not stop another caller from being let through.

## STK-261126-T2BD An operator sets the rate and the burst at start, and a bad setting stops the start

Status: active

The operator sets the sustained rate and the burst when starting
ratelimit, without rebuilding it; a setting it cannot read stops it at the
start with the setting named, rather than limiting by a value nobody
chose.

Why: the limits belong to the deployment, not the build, and a limiter
that silently ran on a default would hide a typo until the service fell
over.

Scenarios:
1. `burst-set` (agent, provisional): an operator who sets the burst to two
   sees the third request at once denied.
2. `bad-setting-stops` (agent, provisional): a rate that is not a number
   stops the start with a failing exit and names the setting.
