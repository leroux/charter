kind: neutral
Feature request (issue #113): select several independent paths at once and join their results into one JSON document.

Given

    {"info":{"friends":[{"first":"Dale","last":"Murphy","age":44},{"first":"Roger","last":"Craig","age":68},{"first":"Jane","last":"Murphy","age":47}]}}

- `[info.friends.0.first,info.friends.1.last]` returns `["Dale","Craig"]`,
- `{info.friends.0.first,info.friends.1.last}` returns `{"first":"Dale","last":"Craig"}` (an object member is named after the last key of its path),
- members can be renamed: `{"alt1":info.friends.0.first,"alt2":info.friends.1.last}` returns `{"alt1":"Dale","alt2":"Craig"}`; a key may also be written without quotes (`{abc:info.friends.0.first}`) and a member whose path has no simple last key (a nested selector, say) is named `_`,
- it combines with any other path component: `info.friends.[0.first,1.age]` returns `["Dale",68]` and `info.friends.#.[first,age]` returns `[["Dale",44],["Roger",68],["Jane",47]]`,
- `[]` and `{}` alone return `[]` and `{}`, and a path that finds nothing is left out of the result,
- a selector may be followed by modifiers, for example `{...}.@pretty.@ugly`.

Please implement it, document it, and add tests.
