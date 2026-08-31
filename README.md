# Gooo evolution trial

This repository runs one bounded, release-to-release Gooo self-improvement
experiment. The loop is declared by `.gooo` metacode; Go is only the executor
and adapter. GitHub Actions is the verification authority.

The experiment observes the released reflexive compiler slice, asks the
released language delta forge for a `NormalizeSource` → `ParseSource` plus
`ValidateStableIDs` split, materializes that candidate only under runner-owned
temporary storage, and sends its exact observations through the released
causal verification runner. It never writes to an upstream checkout.

The result is accepted only when the released compiler can execute the
candidate, replay digests match, the released `CLOSED`/`UNKNOWN`/`REFUTED`
corpus is preserved, causal selection agrees with the full oracle, rollback is
possible, and an exact matched integer resolution pair is observed. A bounded
phase result does not claim whole-language self-improvement or external
utility; those remain `UNKNOWN`/`NOT_MADE`.

Local build and test execution is intentionally outside the development
contract. The workflow runs Go 1.27 and records exact integer measurements.
The root README is excluded from inventory counts.
