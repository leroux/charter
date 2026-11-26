# Modifiers

Code: gjson.go

The `@` components of a path: the built-in modifiers, the ones a caller
adds, and the switch that turns them off.

## SW-261126-H629 A path starting with `@name` runs that modifier and reads the rest of the path from its output

Status: active
Serves: STK-261126-ZQJR

A path, or a path after a pipe, that starts with `@` and a modifier's
name runs the modifier on the JSON so far; an argument may follow the
name after a `:`, as JSON or plain characters. The output is read as a
document of its own by the rest of the path after a `|`. The built-in
modifiers are `@reverse`, which reverses an array's elements or an
object's members and leaves any other value as it was, `@ugly`, which
removes all whitespace, and `@pretty`, which indents, taking the
options `sortKeys`, `indent`, `prefix` and `width` as a JSON object.

Why: a modifier reshapes the JSON where the path syntax alone cannot,
and treating its output as a fresh document lets every other part of
the syntax apply to it.

Scenarios:
1. `reverse` (agent, provisional): `children|@reverse` returns
   `["Jack","Alex","Sara"]`, and `@reverse` on a string returns the
   string.
2. `ugly` (agent, provisional): `@ugly` on an indented document returns
   it with no whitespace between tokens.
3. `pretty` (agent, provisional): `@pretty:{"sortKeys":true}` returns
   the document indented with its keys sorted, and `@ugly` of that is
   the compact document again.

## SW-261126-4Z6X A caller adds a modifier by name, and a name no modifier has gives a result that does not exist

Status: active
Serves: STK-261126-ZQJR

AddModifier registers a function from a JSON text and an argument to a
JSON text under a name, replacing any modifier of that name, and the
name is then usable as `@name` in any path; ModifierExists reports
whether a name is registered. A path naming a modifier nobody
registered gives a result that does not exist.

Why: a caller's own reshaping belongs in its own code, not in gjson,
and an unknown name that quietly returned the input would hide a typo.

Scenarios:
1. `custom` (agent, provisional): a `case` modifier added by the caller
   turns `children|@case:upper` into `["SARA","ALEX","JACK"]`.
2. `exists` (agent, provisional): ModifierExists reports a registered
   name and not an unregistered one.
3. `unknown` (agent, provisional): `@nope` on any document returns a
   result that does not exist.

## SW-261126-Z3CC DisableModifiers makes `@` an ordinary character of a key

Status: active
Serves: STK-261126-ZQJR

With DisableModifiers set, Get runs no modifier and a component
starting with `@` is a key like any other.

Why: a document whose keys start with `@`, as JSON-LD's do, must stay
reachable by a caller who never uses modifiers.

Scenarios:
1. `disabled` (agent, provisional): with DisableModifiers set,
   `@ugly` on `{"@ugly":1}` returns `1`.
