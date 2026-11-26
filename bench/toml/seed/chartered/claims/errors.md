# Errors

Code: errors.go, decode.go

What a refusal carries back: where the problem is, in what form, and the
rule that a refused document gives the program nothing.

## SW-261126-RVMG A refusal names the line and column where the problem begins

Status: active
Kind: bound
Serves: STK-261126-HR5J

Every refusal is a *ParseError holding the line and the column, both
counted from 1, of the first character of what is wrong, and a message
that says what. Columns count characters, not bytes. Its text is
`toml: line L, column C: message`.

Why: a person fixes a configuration file in an editor that shows lines
and columns; a refusal without them leaves the person searching.

Scenarios:
1. `line-and-column` (agent, provisional): a duplicate key on line 3,
   column 1 is refused at line 3, column 1.
2. `columns-count-characters` (agent, provisional): a problem after an é
   on its line is reported one column after it, not two.
3. `message-form` (agent, provisional): the error's text is
   `toml: line L, column C: ` and the message.

## SW-261126-BPDX A refused document gives no table at all

Status: active
Kind: invariant
Serves: STK-261126-HR5J

Decode returns either the whole document or an error and no table: a
document refused on its last line gives the program none of the lines
before it.

Why: a partial configuration is a configuration nobody wrote, and a
program that ran on the first half of a file would fail far from the
cause.

Scenarios:
1. `nothing-on-refusal` (agent, provisional): a document whose last line
   is invalid decodes to a nil table and an error.
