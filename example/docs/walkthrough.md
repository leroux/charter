# How a request spends a token

A walkthrough of ratelimit's code, from a line on standard input to the
answer on standard output. It names the decisions each step carries; the
decisions themselves are in `claims/`.

## Starting

The command reads the rate and the burst from the environment once, in
[`config`](../cmd/ratelimit/main.go#L33), and refuses to start on a value it
cannot read: [SW-261126-CMTG](../claims/config.md#sw-261126-cmtg-the-rate-and-the-burst-come-from-the-environment-at-start-and-an-unreadable-one-stops-the-start). The limiter it builds,
[`New`](../limiter/limiter.go#L28), refuses a rate that is not above zero and a
burst below one.

## One request

Each line of input is a key, and each key is asked of
[`Limiter.Allow`](../limiter/limiter.go#L45). That one method holds the whole
bucket:

<!-- charter:embed SW-261126-VX8Z -->
> **[SW-261126-VX8Z](../claims/limiter.md#sw-261126-vx8z-a-keys-bucket-starts-full-earns-the-rate-holds-at-most-the-burst-and-pays-a-token-a-request) A key's bucket starts full, earns the rate, holds at most the burst and pays a token a request**
>
> A key the limiter has not seen gets a bucket holding the burst. The
> bucket earns the rate in tokens per second of the clock's time, never
> holding more than the burst. A request is let through when the bucket
> holds a whole token and spends it; otherwise it is denied and spends
> nothing.
<!-- /charter:embed -->

It takes the lock for the whole of the request, so one key's bucket is
never another's, under concurrent use too: [SW-261126-39F6](../claims/limiter.md#sw-261126-39f6-one-keys-requests-never-change-another-keys-bucket).

## Time

`Allow` reads the time from the limiter's
[`Clock`](../limiter/clock.go#L9) and nowhere else, which is why a test can
earn a bucket a second's tokens without waiting one:
[SW-261126-2XG8](../claims/limiter.md#sw-261126-2xg8-the-limiter-reads-time-through-its-clock-alone). The command passes the wall clock,
[`SystemClock`](../limiter/clock.go#L14).

<!-- charter:curated
by: agent
date: 2026-11-27
SW-261126-2XG8 20759419089e3372
SW-261126-39F6 6e76dbae670d6652
SW-261126-CMTG f7d5cbb1c616f7ae
SW-261126-VX8Z f5e0ba82ae942577
cmd/ratelimit/main.go#config 49930af804aad8c9
limiter/clock.go#Clock 958fe8350c1998a5
limiter/clock.go#SystemClock d586c09b29898cd6
limiter/limiter.go#Limiter.Allow f25c00572fc8ef71
limiter/limiter.go#New 74d6474104ba7ce6
-->
