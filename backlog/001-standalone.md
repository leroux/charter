# 001. The charter stands alone: a method any repository can take up, with nothing of one product's layout left in it

Status: in progress · Size: L

The tooling builds and passes its own gate here, as `github.com/leroux/charter`. An adopter requires the
module, and `example/` is one (CH-261126-NS7T). What is left:

## Done

- The library at the module root, the charter's own charter in `claims/`,
  the commands as thin mains over `cli`.
- `cli.Claimsgen` and `cli.Claimsmap` take an adopter's charters; a
  charter without a root is read above the working directory (CH-040).
- `example/`: ratelimit, two needs, four decisions (a bound, an invariant,
  a structure and one unkinded), env surfaces accounted, its own
  claimsgen; its gate runs among
  the tooling's tests, and a copy missing a citation must fail it.
- Q-261015-3Q84 answered: the charter is meant for other repositories.

## To do

1. **The method as a document.** The form, the loop and the Style are
   written nowhere in this repository yet. Write them as `docs/method.md`,
   told with the example's entries.
2. **Lint points at the adopter's README.** The style checks name
   `claims/README.md:1` for their counts and "(README, Style)" for the
   bar, a file an adopter need not have. Point them at the method's Style
   section, or let Options name where the bars are written. This changes
   CH entries in `claims/tooling.md` and their tests.
3. **The skills.** `.claude/skills/claims-*` name one product's paths
   (paths, test suites and commands of one layout).
   Make each one read the adopter's layout from its CLAUDE.md or gate
   options, and ship them so an adopter can install them (a Claude Code
   plugin in this repository is one way).
4. **An adopter's CLAUDE.md section.** The "Charter" section an
   adopter puts in its CLAUDE.md, as a template the method document gives.
5. **The policy and the judge in CI.** The workflow's judge reads
   `claims/policy.md`, which the charter's own charter does not have yet.
   Write one, or make the judge's step say there is none.
6. **TypeScript in the example.** The citation extractor is tested on the
   fixture only; a small TypeScript part of the example, citing through a
   generated `claims/index.ts`, would exercise it as an adopter uses it.
7. **The open questions** of `claims/questions.md` (`go run
   ./cmd/claimsgen questions`), several of which were asked before the
   charter was a repository of its own and may now read differently.
8. **Releases.** Tag a version, so an adopter requires one instead of a
   pseudo-version, and say in the README what a version promises.
