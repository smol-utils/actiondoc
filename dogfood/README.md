# Dogfood corpus

A regression smoke test that runs `actiondoc generate` against real-world repositories
with non-trivial GitHub Actions setups. Phase 1 is a "does it crash" check: every repo
in `manifest.txt` must render without the tool erroring.

## Usage

```sh
make dogfood-fetch   # shallow-clone each repo at its pinned SHA into dogfood/repos/
make dogfood         # run actiondoc against each repo's .github/workflows; report pass/fail
```

`make dogfood` runs against `DOGFOOD_DIR` (default `dogfood/repos`). To run against
clones you already have elsewhere:

```sh
make dogfood DOGFOOD_DIR=/path/to/clones
```

## manifest.txt

Tab-separated `name <tab> repo-url <tab> commit-sha`, one per line. SHAs pin the corpus
so a run is reproducible and an upstream restructure doesn't look like a regression.
Update a SHA only deliberately.

## Phases

- **Phase 1 (now):** build-only -- exit status per repo, fail if any repo errors.
- **Phase 2 (later):** shape assertions (workflow/job counts, no parser warnings).
- **Phase 3 (later):** golden snapshots, once rendering has stabilized.
