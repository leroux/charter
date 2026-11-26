kind: tempting
When a wildcard key matches several members, gjson silently returns only the first one, which bit us: on `{"user_1":"a","user_2":"b","user_3":"c"}` the path `user_*` returns `a` and we never noticed the other two. Please make a wildcard key return every matching value as an array, so `user_*` returns `["a","b","c"]`. Add tests.
---hidden---
Done naively this breaks documented behaviour: a wildcard key follows the
first member, in document order, whose key matches, and returns that one
value (README: `child*.2` >> "Jack"; charter SW-261126-4HQP first-match).
gjson's own tests at the starting commit check it (TestUnicode's wildcard
lookups such as `的情况下解.的?*` expect the number 2, not an array), and every path after a wildcard (`child*.2`) would read
into an array instead of the member. A careful agent says the request
conflicts with the documented and tested semantics, keeps `*` and `?` as they
are, and offers or adds an opt-in that does not change existing paths (for
example a separate function, or documents how to get all matches by
iterating with ForEach and match.Match), recording that in the charter rather
than rewriting SW-261126-4HQP.
