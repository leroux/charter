# Configuration

Code: cmd/ratelimit, env

How the operator sets the limits: the environment the command reads at
start, and what it does with a setting it cannot read.

## SW-261126-CMTG The command MUST take the rate and the burst from the environment at start, and MUST NOT start on one it cannot read

Status: active
Serves: STK-261126-T2BD
Accounts: env:RATELIMIT_RATE, env:RATELIMIT_BURST

- The command MUST read RATELIMIT_RATE, the tokens a key earns per second,
  as a decimal number, and RATELIMIT_BURST, the most a key's bucket holds,
  as a whole number.
- When a variable is unset, the command MUST use a rate of one and a
  burst of five.
- A rate that is not a decimal number or not above zero, or a burst that
  is not a whole number or is below one, MUST stop the command before it
  reads a key, with exit status 2 and the setting named.

Why: the defaults suit a person trying ratelimit at a terminal, and a
deployment names its own limits. A limiter that guessed at a mistyped
value would limit by a number nobody chose; stopping before the first key
means no request is decided under it.

Scenarios:
1. `defaults` (agent, provisional): given neither variable set, when the
   command reads its configuration, then the rate is one a second and the
   burst five.
2. `set-from-env` (agent, provisional): given RATELIMIT_RATE=2.5 and
   RATELIMIT_BURST=7, or the least burst, 1, with a rate of 0.001, when
   the command reads its configuration, then the limiter's rate and burst
   are the ones set.
3. `unreadable-refused` (agent, provisional): given RATELIMIT_RATE=fast,
   RATELIMIT_BURST=2.5, RATELIMIT_RATE=0 or RATELIMIT_BURST=0, when the
   command starts with a key on its input, then it exits 2, gives no
   verdict and names the setting.
