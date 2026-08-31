# Gooo evolution trial dossier

This is the committed policy and interpretation for the second immutable
release-to-release experiment. The run-specific machine report and generated
dossier are uploaded by GitHub Actions.

The prior immutable `gooo-evolution-trial v0.1.0` release remains preserved as
a historical `REFUTED` counterexample:
https://github.com/kimjooyoon/gooo-evolution-trial/releases/tag/v0.1.0. Its
exact compiler error was `phase graph must declare exactly three executable
activities`. Evidence disappearance is not closure.

The second experiment uses the immutable `gooo-reflexive-compiler-slice v0.2.0`
release, the exact three-activity baseline phase preserved from the v0.1.1
compiler evidence, the unchanged immutable `gooo-language-delta-forge v0.1.2`
contract, and the unchanged immutable
`gooo-causal-verification-runner v0.1.1` contract. The failure receipt, delta,
and candidate inputs are byte-for-byte the same as the v0.1.0 experiment.

The immutable delta forge emits one exact split candidate: retire the coarse
`NormalizeSource` cell and add two cells, `ParseSource` and
`ValidateStableIDs`. The candidate carries an exact inverse rollback and an
exact integer resolution pair of 1 before and 2 after phase-localization
cells, bound to the same source and toolchain digests.

The candidate `.gooo` phase is generated from the released candidate bundle by
the v0.2.0 released adapter. Candidate Go is never hand-authored as semantic
authority. The v0.2.0 compiler accepts the resulting four-activity phase,
emits semantic IR and backend artifacts, and independently verifies replay.

Both the baseline three-activity corpus and candidate four-activity corpus
must observe exactly one `CLOSED`, one `UNKNOWN`, and one `REFUTED`, with replay
digests matching and no failed commands. The causal runner must report
`total=2`, `selected=1`, `executed=1`, `reused=1`, `full_oracle=2`,
`failures=0`, and `unknowns=0`; `CLOSED` is valid only when the independent full
oracle also reports zero failures and zero unknowns.

The closure receipt is accepted only when observed with stage `IMPROVEMENT`,
step `RESOLVE_TRIAL_COUNTEREXAMPLE`, and reason
`GRAPH_SEMANTICS_ACCEPT_SPLIT_CANDIDATE`. The bounded result does not establish
whole-language self-improvement, which remains `UNKNOWN`; external utility is
`NOT_MADE`.
