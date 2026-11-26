kind: neutral
Modifiers do not work inside the new multi-path selectors. With a custom modifier `case` (upper-cases its input for the argument `upper`) and

    {"friends":[{"age":44,"first":"Dale","last":"Murphy"},{"age":68,"first":"Roger","last":"Craig"},{"age":47,"first":"Jane","last":"Murphy"}]}

- `friends.#.{age,first|@case:upper}|@ugly` should return `[{"age":44,"@case:upper":"DALE"},{"age":68,"@case:upper":"ROGER"},{"age":47,"@case:upper":"JANE"}]` (a member is named after the last component of its path, here the modifier),
- `{friends.#.{age,first:first|@case:upper}|0.first}` should return `{"first":"DALE"}`.

Please fix and add tests.
