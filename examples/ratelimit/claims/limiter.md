# The limiter

Code: limiter

How the limiter decides, per key, whether a request goes through: the
token bucket each key has, and the clock it reads time from.

## SW-261126-VX8Z The limiter MUST give a new key a full bucket that earns the rate, holds at most the burst and pays a token for each request let through

Status: active
Kind: bound
Serves: STK-261126-H2SW

- The limiter MUST give a key it has not seen a bucket holding the burst.
- The bucket MUST earn the rate in tokens per second of the clock's time,
  and MUST NOT hold more than the burst.
- The limiter MUST let a request through when the key's bucket holds a
  whole token, and spend that token.
- The limiter MUST deny a request when the bucket holds less than a whole
  token, and MUST NOT spend anything for it.

Why: a bucket that starts full lets a caller that has been quiet burst at
once. Capping the bucket at the burst stops a caller that was quiet for an
hour from saving up an hour of requests and sending them together. A
denial spends nothing, so a caller that asks too early is not pushed
further back for having asked.

Scenarios:
1. `starts-full` (agent, provisional): given a burst of three and a new
   key, when four requests for it arrive at one instant, then three are
   let through and the fourth is denied.
2. `earns-the-rate` (agent, provisional): given a rate of two, a burst of
   three and a key whose bucket is empty, when one second passes and
   three requests arrive, then two are let through and the third is
   denied.
3. `capped-at-burst` (agent, provisional): given a rate of ten, a burst
   of four and a key whose bucket is empty, when the key is quiet for an
   hour and then sends a hundred requests at once, then four are let
   through.
4. `denial-spends-nothing` (agent, provisional): given a rate of one, a
   burst of one and an empty bucket, when a request half a second later
   is denied and another half second passes, then the next request is let
   through.

## SW-261126-39F6 The limiter MUST NOT let one key's requests change another key's bucket

Status: active
Kind: invariant
Serves: STK-261126-H2SW

Spending, earning and denial on one key MUST leave every other key's
bucket as it was, and MUST do so when many keys are used at once.

Why: a limiter that let keys share tokens would deny a caller because a
neighbour had spent them, which is what a per-caller limit exists to
prevent.

Scenarios:
1. `keys-independent` (agent, provisional): given a burst of two and a
   key that has sent ten requests, when a new key sends three at the same
   instant, then two of the new key's are let through.
2. `concurrent-keys` (agent, provisional): given a burst of five and
   eight keys, when four concurrent callers per key each send five
   requests at once, then exactly five of each key's twenty are let
   through.

## SW-261126-2XG8 The limiter MUST read the time from the Clock it is given and nowhere else

Status: active
Kind: structure
Serves: STK-261126-H2SW

The limiter MUST take a Clock when it is made and MUST read the time from
that Clock and nowhere else. The command gives it the wall clock.

Why: a limit is a statement about time, and a test that had to sleep to
check one would be slow and flaky. A limiter that reads only its Clock
lets a test move the time by hand, so every bound can be checked at an
instant.

Scenarios:
1. `time-from-clock` (agent, provisional): given a rate of a thousand, a
   burst of one and a Clock that does not move, when the bucket is
   emptied and twenty milliseconds of wall time pass, then the next
   request is denied.
