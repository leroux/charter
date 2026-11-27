# 006. Ideas held for later: extensions the foundations do not need yet

Status: held · Size: — · See: 007, 008

Each idea came out of the design talk of 2026-11-26. None is ruled; each
waits for an adopter, torii first, to show it is needed. Taking one up is
a story of its own, drafted from what that adopter showed.

| Idea | What it would be | What would show it is needed | Where it was argued |
|---|---|---|---|
| Contract and design | An area says `Level: contract` or `Level: design`; surfaces belong to contract entries | A record with real architecture choices to keep apart from what binds every implementation, as torii's may be | design talk; drafted as an entry and deleted unbuilt |
| Portable fixtures | Contract scenarios as txtar files of a tree, a command and the expected output, run by `charter conform` against any implementation | A reimplementation to check | design talk; drafted as an entry and deleted unbuilt |
| The edge declared complete | The charter's own verbs, flags, charter.toml keys and JSON shapes as surface kinds | Wanting the gate to name behaviour no entry describes | design talk |
| Configurable levels | `[[level]]` in charter.toml: name, file, anchor, how it is verified, whether Serves is required; outcome and decision by default | A project whose levels do not fit two, such as one verified by an evaluation | design talk |
| Composed charters | A charter per component, a parent listing them, and a checked link to a child's outcome | A system too large, or too widely owned, for one record with areas | design talk; observations may serve instead |
| Nested areas and checked Code lines | Area documents in folders, and an anchor outside its area's Code paths named | A record large enough that areas need grouping, or Code lines found drifting again | design talk |
| A machine interface | `charter check --json`, a JSON export of the charter, and facts taken in, for other tools and forks of spec kits | Another tool, or a fork of OpenSpec or Spec Kit, building on the charter | design talk |
| A bridge to OpenSpec | A reader of OpenSpec's living specs as a record | Others wanting charter's gate over specs they already keep | design talk |
| Change proposals | A change as a folder: why, the entries added, modified and removed, and its tasks | The backlog, Changes lines and diff proving too loose for many people at once | design talk |
| Outcomes in place of needs | Needs renamed outcomes, worded with no person, verified from edge-tests | Settling the record's top layer | design talk |
