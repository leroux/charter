# How a request spends a token

A walkthrough of ratelimit's code, from a line on standard input to the
answer on standard output. Each step links the decision it carries, and
`claims/` holds the decisions in full.

## Starting

The command reads the rate and the burst from the environment once, in
[`config`](../cmd/ratelimit/main.go#L33), and refuses to start on a value it
cannot read: [SW-261126-CMTG](../claims/config.md#sw-261126-cmtg-the-command-must-take-the-rate-and-the-burst-from-the-environment-at-start-and-must-not-start-on-one-it-cannot-read). The limiter it builds,
[`New`](../limiter/limiter.go#L28), refuses a rate that is not above zero and a
burst below one.

## One request

Each line of input is a key, and each key is asked of
[`Limiter.Allow`](../limiter/limiter.go#L45). That one method carries the whole
bucket rule:

<!-- charter:embed SW-261126-VX8Z -->
> **[SW-261126-VX8Z](../claims/limiter.md#sw-261126-vx8z-the-limiter-must-give-a-new-key-a-full-bucket-that-earns-the-rate-holds-at-most-the-burst-and-pays-a-token-for-each-request-let-through) The limiter MUST give a new key a full bucket that earns the rate, holds at most the burst and pays a token for each request let through**
>
> - The limiter MUST give a key it has not seen a bucket holding the burst.
> - The bucket MUST earn the rate in tokens per second of the clock's time,
>   and MUST NOT hold more than the burst.
> - The limiter MUST let a request through when the key's bucket holds a
>   whole token, and spend that token.
> - The limiter MUST deny a request when the bucket holds less than a whole
>   token, and MUST NOT spend anything for it.
<!-- /charter:embed -->

`Allow` holds the limiter's lock for the whole request, so one key's
requests never change another key's bucket, even when many keys are used
at once: [SW-261126-39F6](../claims/limiter.md#sw-261126-39f6-the-limiter-must-not-let-one-keys-requests-change-another-keys-bucket).

## Time

`Allow` reads the time from the limiter's
[`Clock`](../limiter/clock.go#L9) and nowhere else, which is why a test can
earn a bucket a second's tokens without waiting one:
[SW-261126-2XG8](../claims/limiter.md#sw-261126-2xg8-the-limiter-must-read-the-time-from-the-clock-it-is-given-and-nowhere-else). The command passes the wall clock,
[`SystemClock`](../limiter/clock.go#L14).

<!-- charter:curated
by: agent
date: 2026-11-28
SW-261126-2XG8 53ffd96022a7ba94
SW-261126-39F6 3f2938df147192da
SW-261126-CMTG 24fb8eacccd57cc4
SW-261126-VX8Z c7c40f6af486a26c
cmd/ratelimit/main.go#config 49930af804aad8c9
limiter/clock.go#Clock 958fe8350c1998a5
limiter/clock.go#SystemClock d586c09b29898cd6
limiter/limiter.go#Limiter.Allow f25c00572fc8ef71
limiter/limiter.go#New 74d6474104ba7ce6
-->
