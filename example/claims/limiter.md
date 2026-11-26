# The limiter

Code: limiter

How the limiter decides, per key, whether a request goes through: the
token bucket each key has, and the clock it reads time from.

## SW-261126-VX8Z A key's bucket starts full, earns the rate, holds at most the burst and pays a token a request

Status: active
Kind: bound
Serves: STK-261126-H2SW

A key the limiter has not seen gets a bucket holding the burst. The
bucket earns the rate in tokens per second of the clock's time, never
holding more than the burst. A request is let through when the bucket
holds a whole token and spends it; otherwise it is denied and spends
nothing.

Why: a bucket that starts full lets a quiet caller burst at once, and
one capped at the burst keeps a caller quiet for an hour from saving up
an hour of requests.

Scenarios:
1. `starts-full` (agent, provisional): a new key is let through the
   burst's number of times at one instant and denied the next.
2. `earns-the-rate` (agent, provisional): after an empty bucket waits
   one second at a rate of two, two requests go through and the third is
   denied.
3. `capped-at-burst` (agent, provisional): a key quiet for far longer
   than the burst takes to earn is let through the burst and no more.
4. `denial-spends-nothing` (agent, provisional): a denied request does
   not delay the next token.

## SW-261126-39F6 One key's requests never change another key's bucket

Status: active
Kind: invariant
Serves: STK-261126-H2SW

Each key's bucket is its own: spending, earning and denial on one key
leave every other key's bucket as it was, under concurrent use too.

Why: a limiter that let keys share tokens would punish a caller for a
neighbour, the thing a per-caller limit exists to prevent.

Scenarios:
1. `keys-independent` (agent, provisional): a key that has emptied its
   bucket leaves a new key the full burst.
2. `concurrent-keys` (agent, provisional): many keys requesting at once
   are each let through exactly the burst.

## SW-261126-2XG8 The limiter reads time through its Clock alone

Status: active
Kind: structure
Serves: STK-261126-H2SW

The limiter takes a Clock when it is made and reads the time from it and
nowhere else; the command passes the wall clock.

Why: a limit is a statement about time, and a test that must sleep to
check it is slow and flaky; one clock to drive makes every bound
testable at an instant.

Scenarios:
1. `time-from-clock` (agent, provisional): a limiter given a clock that
   does not move earns nothing, however long the test takes.
