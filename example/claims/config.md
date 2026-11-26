# Configuration

Code: cmd/ratelimit, env

How the operator sets the limits: the environment the command reads at
start, and what it does with a setting it cannot read.

## SW-261126-CMTG The rate and the burst come from the environment at start, and an unreadable one stops the start

Status: active
Serves: STK-261126-T2BD
Accounts: env:RATELIMIT_RATE, env:RATELIMIT_BURST

RATELIMIT_RATE is the tokens a key earns per second, a decimal number,
and RATELIMIT_BURST the most a key holds, a whole number; unset, they are
one and five. A value that is not a number of its form, a rate not above
zero or a burst below one stops the command before it reads a key, with
exit status 2 and the setting named.

Why: the defaults suit a person trying it at a terminal, and a deployment
names its own; a limiter that guessed at a mistyped value would limit by
a number nobody chose.

Scenarios:
1. `defaults` (agent, provisional): with neither set, the limit is one a
   second with a burst of five.
2. `set-from-env` (agent, provisional): a rate and a burst set in the
   environment are the limiter's, the least of each, any rate above zero
   and a burst of one, among them.
3. `unreadable-refused` (agent, provisional): a rate not a number, a
   burst not whole, and a rate or burst of zero each exit 2 before any
   key is read, naming the setting.
