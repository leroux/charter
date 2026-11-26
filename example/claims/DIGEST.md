<!-- Generated from the charter's documents. DO NOT EDIT. Regenerate: charter gen -->

# The charter, read through

What the system does, without its code: 2 needs, what a person or the operator can rely on, and 4 decisions in 2 areas, how the system behaves at each boundary. Read the needs, then the map of the areas, then the areas in order. Each decision is its title, the sentence it rules, and its why; its id links to the whole entry, with its scenarios.

1. [What people need](#what-people-need)
2. [The areas](#the-areas)
3. [The decisions, area by area](#the-decisions-area-by-area)

## What people need

Each need is listed under the area that does most to meet it.

**The limiter**

- [STK-261126-H2SW](needs.md#stk-261126-h2sw-a-caller-gets-a-burst-at-once-then-its-sustained-rate-and-no-more) A caller gets a burst at once, then its sustained rate, and no more

**Configuration**

- [STK-261126-T2BD](needs.md#stk-261126-t2bd-an-operator-sets-the-rate-and-the-burst-at-start-and-a-bad-setting-stops-the-start) An operator sets the rate and the burst at start, and a bad setting stops the start

## The areas

| # | Area | What it holds | Decisions |
|---|---|---|---|
| 1 | [The limiter](limiter.md) | How the limiter decides, per key, whether a request goes through: the token bucket each key has, and the clock it reads time from. | 3 |
| 2 | [Configuration](config.md) | How the operator sets the limits: the environment the command reads at start, and what it does with a setting it cannot read. | 1 |

Which areas serve each need, with how many of their decisions:

| Need | Areas |
|---|---|
| [STK-261126-H2SW](needs.md#stk-261126-h2sw-a-caller-gets-a-burst-at-once-then-its-sustained-rate-and-no-more) | limiter 3 |
| [STK-261126-T2BD](needs.md#stk-261126-t2bd-an-operator-sets-the-rate-and-the-burst-at-start-and-a-bad-setting-stops-the-start) | config 1 |

## The decisions, area by area

### 1. The limiter

How the limiter decides, per key, whether a request goes through: the token bucket each key has, and the clock it reads time from.

- **[SW-261126-VX8Z](limiter.md#sw-261126-vx8z-a-keys-bucket-starts-full-earns-the-rate-holds-at-most-the-burst-and-pays-a-token-a-request)** A key's bucket starts full, earns the rate, holds at most the burst and pays a token a request *(bound)*  
  Why: a bucket that starts full lets a quiet caller burst at once, and one capped at the burst keeps a caller quiet for an hour from saving up an hour of requests.
- **[SW-261126-39F6](limiter.md#sw-261126-39f6-one-keys-requests-never-change-another-keys-bucket)** One key's requests never change another key's bucket *(invariant)*  
  Why: a limiter that let keys share tokens would punish a caller for a neighbour, the thing a per-caller limit exists to prevent.
- **[SW-261126-2XG8](limiter.md#sw-261126-2xg8-the-limiter-reads-time-through-its-clock-alone)** The limiter reads time through its Clock alone *(structure)*  
  Why: a limit is a statement about time, and a test that must sleep to check it is slow and flaky; one clock to drive makes every bound testable at an instant.

### 2. Configuration

How the operator sets the limits: the environment the command reads at start, and what it does with a setting it cannot read.

- **[SW-261126-CMTG](config.md#sw-261126-cmtg-the-rate-and-the-burst-come-from-the-environment-at-start-and-an-unreadable-one-stops-the-start)** The rate and the burst come from the environment at start, and an unreadable one stops the start  
  Why: the defaults suit a person trying it at a terminal, and a deployment names its own; a limiter that guessed at a mistyped value would limit by a number nobody chose.
