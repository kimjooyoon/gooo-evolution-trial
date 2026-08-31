#!/usr/bin/env bash
set -Eeuo pipefail

root=${1:?repository root is required}
trial=${2:?trial executable is required}
work=${3:?experiment work directory is required}
mkdir -p "$work/bin" "$work/src" "$work/results" "$work/logs"

bash "$root/scripts/download-upstreams.sh" "$root" "$work/upstream"

tar -xzf "$work/upstream/assets/gooo-language-delta-forge-v0.1.2.tar.gz" -C "$work/src"
tar -xzf "$work/upstream/assets/gooo-reflexive-compiler-slice-source-v0.1.1.tar.gz" -C "$work/src"
tar -xzf "$work/upstream/assets/gooo-causal-verification-runner-v0.1.1.tar.gz" -C "$work/src"

delta_root="$work/src/gooo-language-delta-forge-v0.1.2"
compiler_root="$work/src/gooo-reflexive-compiler-slice-v0.1.1"
causal_root="$work/src/gooo-causal-verification-runner-v0.1.1"

(cd "$delta_root" && go build -o "$work/bin/delta-forge" ./cmd/gooo-language-delta-forge)
(cd "$compiler_root" && go build -o "$work/bin/reflexive-compiler" ./cmd/gooo-reflexive-compiler-slice && go build -o "$work/bin/reflexive-verify" ./cmd/gooo-reflexive-verify)
(cd "$causal_root" && go build -o "$work/bin/causal-runner" ./cmd/gooo-causal-verification-runner)

"$trial" authority-verify --meta "$root/meta/evolution-trial.gooo" --contract "$root/contracts/evolution-trial-denominator-v1.json" --output "$work/results/authority-report.json" > "$work/logs/authority.stdout"
"$trial" inventory --root "$root" > "$work/results/repository-inventory.json"

mkdir -p "$work/results/delta-conformance"
"$work/bin/delta-forge" conformance \
	--program "$delta_root/examples/language-delta-forge-v1/main.gooo" \
	--denominator "$delta_root/contracts/denominator-v1.json" \
	--manifest "$delta_root/fixtures/conformance.json" \
	--output "$work/results/delta-conformance" > "$work/logs/delta-conformance.stdout"

mkdir -p "$work/results/causal-conformance"
"$work/bin/causal-runner" conformance \
	--source "$causal_root/examples/causal-verification/main.gooo" \
	--contract "$causal_root/contracts/causal-verification-denominator-v1.json" \
	--corpus "$causal_root/fixtures/corpus.json" \
	--tree-root "$causal_root" \
	--output-dir "$work/results/causal-conformance" > "$work/logs/causal-conformance.stdout"

mkdir -p "$work/results/baseline-corpus"
"$trial" run-corpus \
	--compiler "$work/bin/reflexive-compiler" \
	--verifier "$work/bin/reflexive-verify" \
	--phase "$compiler_root/meta/reflexive-normalize.gooo" \
	--root "$compiler_root" --role baseline \
	--output-dir "$work/results/baseline-corpus" \
	--report "$work/results/baseline-report.json" > "$work/logs/baseline-corpus.stdout"

mkdir -p "$work/results/delta-input"
"$trial" prepare-input \
	--phase "$compiler_root/meta/reflexive-normalize.gooo" \
	--compiler-commit "dabbe38badebefdf2979d8862c26a647b0dd15c0" \
	--baseline-receipt "$work/results/baseline-corpus/CLOSED_CANONICAL_INPUT/baseline/receipt.json" \
	--output-dir "$work/results/delta-input" > "$work/logs/prepare-input.stdout"

mkdir -p "$work/results/delta-candidate"
"$work/bin/delta-forge" generate \
	--program "$delta_root/examples/language-delta-forge-v1/main.gooo" \
	--denominator "$delta_root/contracts/denominator-v1.json" \
	--input "$work/results/delta-input/input.json" \
	--output "$work/results/delta-candidate" > "$work/logs/delta-generate.stdout"

"$trial" generate-candidate-phase \
	--candidate "$work/results/delta-candidate/candidate-bundle.json" \
	--output "$work/results/candidate-phase.gooo" > "$work/logs/candidate-phase.stdout"

mkdir -p "$work/results/candidate-corpus"
"$trial" run-corpus \
	--compiler "$work/bin/reflexive-compiler" \
	--verifier /bin/true \
	--phase "$work/results/candidate-phase.gooo" \
	--root "$compiler_root" --role candidate \
	--output-dir "$work/results/candidate-corpus" \
	--report "$work/results/candidate-report.json" > "$work/logs/candidate-corpus.stdout"

"$trial" prepare-causal-case \
	--baseline-report "$work/results/baseline-report.json" \
	--candidate-report "$work/results/candidate-report.json" \
	--causal-release-id 380048457 \
	--causal-asset-name "gooo-causal-verification-runner-v0.1.1.tar.gz" \
	--causal-asset-digest "sha256:fcf40acd1f09805e526b8e9634cd70b8f308e1f2312b0aac8d3253c4038db7fb" \
	--output "$work/results/causal-case.json" > "$work/logs/prepare-causal-case.stdout"

mkdir -p "$work/results/causal-experiment"
"$work/bin/causal-runner" run \
	--source "$causal_root/examples/causal-verification/main.gooo" \
	--contract "$causal_root/contracts/causal-verification-denominator-v1.json" \
	--case "$work/results/causal-case.json" \
	--output-dir "$work/results/causal-experiment" \
	--tree-root "$causal_root" \
	--compile-ms 0 --build-ms 0 --test-ms 0 > "$work/logs/causal-experiment.stdout"

printf '%s\n' '{"schema":"gooo/evolution-trial/released-tools-executed/v1","delta_forge":true,"reflexive_compiler":true,"causal_runner":true,"candidate_only":true,"upstream_writes":0,"repository_writes":0}' > "$work/results/tool-execution-report.json"
