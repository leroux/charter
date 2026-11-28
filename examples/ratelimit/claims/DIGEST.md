<!-- Generated from the charter's documents. DO NOT EDIT. Regenerate: charter gen -->

# The charter, read through

What the system does, without its code: 2 needs, what a person or the operator can rely on, and 4 decisions in 2 areas, how the system behaves at each boundary. Read the needs, then the map of the areas, then the areas in order. Each decision is its title, the sentence it rules, and its why; its id links to the whole entry, with its scenarios.

1. [What people need](#what-people-need)
2. [The areas](#the-areas)
3. [The decisions, area by area](#the-decisions-area-by-area)

## What people need

Each need is listed under the area that does most to meet it.

**The limiter**

- [STK-261126-H2SW](needs.md#stk-261126-h2sw-ratelimit-must-let-a-quiet-caller-burst-then-hold-it-to-its-sustained-rate-and-must-not-let-one-caller-spend-anothers-requests) ratelimit MUST let a quiet caller burst, then hold it to its sustained rate, and MUST NOT let one caller spend another's requests

**Configuration**

- [STK-261126-T2BD](needs.md#stk-261126-t2bd-ratelimit-must-take-its-rate-and-burst-from-the-operator-at-start-and-must-not-start-on-a-setting-it-cannot-read) ratelimit MUST take its rate and burst from the operator at start, and MUST NOT start on a setting it cannot read

## The areas

| # | Area | What it holds | Decisions |
|---|---|---|---|
| 1 | [The limiter](limiter.md) | How the limiter decides, per key, whether a request goes through: the token bucket each key has, and the clock it reads time from. | 3 |
| 2 | [Configuration](config.md) | How the operator sets the limits: the environment the command reads at start, and what it does with a setting it cannot read. | 1 |

Which areas serve each need, with how many of their decisions:

| Need | Areas |
|---|---|
| [STK-261126-H2SW](needs.md#stk-261126-h2sw-ratelimit-must-let-a-quiet-caller-burst-then-hold-it-to-its-sustained-rate-and-must-not-let-one-caller-spend-anothers-requests) | limiter 3 |
| [STK-261126-T2BD](needs.md#stk-261126-t2bd-ratelimit-must-take-its-rate-and-burst-from-the-operator-at-start-and-must-not-start-on-a-setting-it-cannot-read) | config 1 |

## The decisions, area by area

### 1. The limiter

How the limiter decides, per key, whether a request goes through: the token bucket each key has, and the clock it reads time from.

- **[SW-261126-VX8Z](limiter.md#sw-261126-vx8z-the-limiter-must-give-a-new-key-a-full-bucket-that-earns-the-rate-holds-at-most-the-burst-and-pays-a-token-for-each-request-let-through)** The limiter MUST give a new key a full bucket that earns the rate, holds at most the burst and pays a token for each request let through *(bound)*  
  Why: a bucket that starts full lets a caller that has been quiet burst at once.
- **[SW-261126-39F6](limiter.md#sw-261126-39f6-the-limiter-must-not-let-one-keys-requests-change-another-keys-bucket)** The limiter MUST NOT let one key's requests change another key's bucket *(invariant)*  
  Why: a limiter that let keys share tokens would deny a caller because a neighbour had spent them, which is what a per-caller limit exists to prevent.
- **[SW-261126-2XG8](limiter.md#sw-261126-2xg8-the-limiter-must-read-the-time-from-the-clock-it-is-given-and-nowhere-else)** The limiter MUST read the time from the Clock it is given and nowhere else *(structure)*  
  Why: a limit is a statement about time, and a test that had to sleep to check one would be slow and flaky.

### 2. Configuration

How the operator sets the limits: the environment the command reads at start, and what it does with a setting it cannot read.

- **[SW-261126-CMTG](config.md#sw-261126-cmtg-the-command-must-take-the-rate-and-the-burst-from-the-environment-at-start-and-must-not-start-on-one-it-cannot-read)** The command MUST take the rate and the burst from the environment at start, and MUST NOT start on one it cannot read  
  Why: the defaults suit a person trying ratelimit at a terminal, and a deployment names its own limits.
