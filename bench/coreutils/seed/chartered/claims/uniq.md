# uniq

Code: cmd/uniq

How uniq merges runs of repeated lines, and how -c counts them.

## SW-261126-4T6M uniq writes each run of adjacent byte-identical lines once

Status: active
Kind: invariant
Serves: STK-261126-69ZV

Lines that are the same, byte for byte, and adjacent form a run, and uniq
writes the run's line once. Equal lines with another line between them
are separate runs. Lines that differ only in case differ.

Why: GNU's uniq is a filter over adjacent lines, which is why scripts
put sort in front of it; and case matters in the C locale, so "Apple"
and "apple" are two records.

Scenarios:
1. `adjacent-merged` (agent, provisional): a, a, b prints a, b.
2. `apart-stay` (agent, provisional): a, b, a prints a, b, a.
3. `case-matters` (agent, provisional): Apple, apple prints both.

## SW-261126-8785 uniq -c prefixes each line with its run's length in seven columns and a space

Status: active
Serves: STK-261126-69ZV

With -c, each line uniq writes is preceded by the length of its run,
right-aligned in seven columns, and a space.

Why: GNU's format; `sort | uniq -c | sort -n` is among the most common
pipelines there are, and it relies on the count leading the line.

Scenarios:
1. `count-seven-columns` (agent, provisional): a, a, b with -c prints
   "      2 a" and "      1 b".
