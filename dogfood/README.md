# Dogfood corpus

A regression test that runs `actiondoc generate` against real-world repositories with
non-trivial GitHub Actions setups, at two levels:

1. **Crash check** (`make dogfood`): every repo in `manifest.txt` must render without
   the tool erroring.
2. **Output snapshots** (`make dogfood-output`): for a subset of repos chosen to exercise
   every rendering surface, the rendered Markdown is committed under `snapshots/` and
   diffed on every run. A rendering change -- intentional or not -- shows up as a
   reviewable diff.

## Usage

```sh
make dogfood-fetch     # shallow-clone each repo at its pinned SHA into dogfood/repos/
make dogfood           # crash check: render every repo, report pass/fail
make dogfood-output    # snapshot check: diff rendered output against snapshots/
```

All targets run against `DOGFOOD_DIR` (default `dogfood/repos`). To run against clones
you already have elsewhere:

```sh
make dogfood DOGFOOD_DIR=/path/to/clones
```

After an intentional rendering change, regenerate and commit the snapshots:

```sh
make dogfood-output-update
git add dogfood/snapshots && git commit
```

The snapshot diff in that commit is the review artifact: it shows exactly what the
change did to real-world output.

## manifest.txt

Tab-separated `name <tab> repo-url <tab> commit-sha`, one per line. SHAs pin the corpus
so a run is reproducible and an upstream restructure doesn't look like a regression.
Update a SHA only deliberately -- snapshots are only meaningful against the pinned SHAs.

## Snapshot repos

The snapshot subset (see `SNAPSHOT_REPOS` in the Makefile) is chosen so that together the
repos cover every rendering surface where output bugs have been found in the wild:
secrets inventories, multi-line step names, YAML anchors, matrix job names, deep
reusable-workflow call graphs, composite action docs, and link/anchor escaping. Keep the
set small enough that a snapshot diff stays reviewable; extend it when a new bug class is
found in a repo outside the set.
