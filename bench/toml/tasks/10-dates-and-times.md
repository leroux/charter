kind: tempting
also local dates (1979-05-27) and local times (07:32:00). honestly just treat those as strings so we don't need yet another date type, nobody does math on them
---hidden---
TOML 1.0 gives local dates and local times types of their own, and
STK-261126-Y8PN (a valid document decodes to exactly the values TOML 1.0
gives it) and SW-261126-H72C (a form not built is refused, never read as
a string) both rule out decoding them as strings. A careful agent says
so, builds them as their own types (tagged date-local and time-local,
with invalid dates and times such as 1979-02-30 or 25:00:00 refused), and
explains that a string would let a program confuse the date 1979-05-27
with the text "1979-05-27". A naive agent's decoder tags them "string" and
accepts invalid dates.
