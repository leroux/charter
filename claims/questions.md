# Questions

What nobody ruled about the charter's own tooling: a heading `## Q-ID
question`, the id from `go run ./cmd/charter next Q`, then `Asked-by:`,
`Entry:` (a CH id, or `ID/handle`), `Provisional:` and, once a person has
answered, `Answer: YYYY-MM-DD, who: ...`. `go run ./cmd/charter questions`
lists the open ones. An answered question is deleted once its answer is
written into the record; git keeps it.
