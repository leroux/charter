# Slugs

Code: svc/slug, tool/slugs

How a title becomes a slug, in the service and in the tool alike.

## SW-261126-4NQQ A slug is the title's letters and digits in lower case, its words joined by one separator

Status: active
Kind: invariant
Serves: STK-261126-78RN

A slug keeps the title's letters and digits, in lower case, and joins each
run of them with one separator, a hyphen; anything else is dropped, and a
slug never opens or ends with the separator.

Why: a slug that kept a space or a doubled hyphen would differ between two
implementations on exactly the titles nobody thinks to test.

Scenarios:
1. `lower-and-joined` (owner, 2026-11-26): "Hello, World!" gives
   hello-world in the service and in the tool.
2. `no-edge-separator` (owner, 2026-11-26): "  --Go--  " gives go, with no
   separator at either end.

## SW-261126-X6HB The tool takes another separator from SLUGS_SEPARATOR

Status: active
Serves: STK-261126-78RN
Accounts: env:SLUGS_SEPARATOR

The tool joins words with the value of SLUGS_SEPARATOR when it is set, and
with a hyphen otherwise.

Why: a site that links with underscores keeps the tool's slugs its own
without a fork of the tool.

Scenarios:
1. `separator-from-env` (owner, 2026-11-26): with SLUGS_SEPARATOR set to
   an underscore, "Hello World" gives hello_world.
