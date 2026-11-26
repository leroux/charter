# Needs

What the author of a shell script that calls cut, sort, uniq or wc can
rely on. These tools are drop-in clones of GNU coreutils 9.4 run in the C
locale (LC_ALL=C): a script written against GNU must not notice the swap.
A need is verified from the script's side: a tool's command, given its
arguments and input and read back.

## STK-261126-69ZV A script gets GNU's output, byte for byte, for every flag a tool supports

Status: active

For every flag a tool supports, its standard output and exit status are
GNU's for the same arguments and input: the same bytes, the same
padding, the same order. A tool supports fewer flags than GNU, never a
different meaning for one.

Why: scripts parse this output with other tools, and a column that moved
by one space or a line in another order breaks them silently; the value
of a clone is that nobody has to read it to trust it.

Scenarios:
1. `wc-like-gnu` (agent, provisional): wc with no flag prints a file's
   lines, words and bytes as GNU prints them.
2. `uniq-count-like-gnu` (agent, provisional): uniq -c over sorted lines
   gives the bytes GNU's gives.
3. `cut-like-gnu` (agent, provisional): cut -d: -f1,3 gives the bytes
   GNU's gives, lines without a colon among them.

## STK-261126-3NDQ A tool reads standard input when no file is named, the named files in order otherwise, and loses no line

Status: active

With no file named a tool reads standard input; with files named it reads
each in the order given. Every line of the input reaches the output it
belongs in, the last one included when it lacks its newline.

Why: a script feeds these tools both ways, from a pipe and from files,
and a lost last line is a lost record nobody sees go.

Scenarios:
1. `stdin-when-no-file` (agent, provisional): a tool with no file named
   reads its standard input.
2. `files-in-order` (agent, provisional): cut over two files prints the
   first file's lines, then the second's.
3. `last-line-kept` (agent, provisional): a last line without a newline
   is printed, with a newline.

## STK-261126-WSGT A flag a tool does not support stops it, rather than printing output nobody asked for

Status: active

Given an option it does not know, a tool writes nothing to standard
output, names the option on standard error and exits with a failing
status.

Why: a script that passes a flag relies on its meaning; output computed
without it would be wrong and look right.

Scenarios:
1. `unknown-flag-fails` (agent, provisional): each tool given -Z exits
   non-zero, prints nothing on standard output and names Z on standard
   error.
