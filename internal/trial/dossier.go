package trial

import (
	"fmt"
	"os"
	"strings"
)

func WriteDossier(baselinePath, candidatePath, candidateBundlePath, evidencePath, causalReceiptPath, processEvidencePath, metricsPath, reportPath, dossierPath string) error {
	var baseline CorpusReport
	var candidate CorpusReport
	var bundle CandidateBundle
	if err := ReadJSON(baselinePath, &baseline); err != nil {
		return err
	}
	if err := ReadJSON(candidatePath, &candidate); err != nil {
		return err
	}
	if err := ReadJSON(candidateBundlePath, &bundle); err != nil {
		return err
	}
	var evidence map[string]any
	if err := ReadJSON(evidencePath, &evidence); err != nil {
		return err
	}
	var causal map[string]any
	if err := ReadJSON(causalReceiptPath, &causal); err != nil {
		return err
	}
	var process map[string]any
	if err := ReadJSON(processEvidencePath, &process); err != nil {
		return err
	}
	var metrics map[string]any
	if err := ReadJSON(metricsPath, &metrics); err != nil {
		return err
	}
	causalDecision, _ := causal["decision"].(string)
	causalOracleState := "UNKNOWN"
	if comparison, ok := causal["full_oracle_comparison"].(map[string]any); ok {
		if state, ok := comparison["state"].(string); ok {
			causalOracleState = state
		}
	}
	resolutionObserved := false
	if pair, ok := evidence["resolution_pair"].(map[string]any); ok {
		_, beforeOK := pair["before"]
		_, afterOK := pair["after"]
		resolutionObserved = beforeOK && afterOK
	}
	acceptance := map[string]any{
		"candidate_compiled_and_generated_ir_backend": candidate.InterfaceDecision == "CLOSED" && allGenerated(candidate),
		"replay_digests_match":                        baseline.ReplayDigestMatch && candidate.ReplayDigestMatch,
		"closed_unknown_refuted_corpus_preserved":     baseline.Closed == 1 && baseline.Unknown == 1 && baseline.Refuted == 1 && candidate.Closed == 1 && candidate.Unknown == 1 && candidate.Refuted == 1,
		"causal_selection_and_full_oracle_closed":     causalDecision == "CLOSED" && causalOracleState == "CLOSED",
		"rollback_possible":                           bundle.RollbackDelta.ExactPair,
		"exact_semantic_resolution_pair_observed":     resolutionObserved,
	}
	decision := "CLOSED"
	if candidate.InterfaceDecision == "REFUTED" || candidate.Failed > 0 || causalDecision == "REFUTED" {
		decision = "REFUTED"
	} else if candidate.InterfaceDecision == "UNKNOWN" || causalDecision == "UNKNOWN" {
		decision = "UNKNOWN"
	}
	final := map[string]any{
		"schema":                   "gooo/evolution-trial/final-report/v1",
		"decision":                 decision,
		"precedence":               []string{"REFUTED", "UNKNOWN", "CLOSED"},
		"experiment":               "first-release-to-release-reflexive-normalization-split",
		"authority":                map[string]any{"repository_writes": 0, "upstream_writes": 0, "local_test_executions": 0, "output_location": "CALLER_OWNED_TEMP_ONLY", "verification_authority": "GITHUB_ACTIONS"},
		"candidate":                map[string]any{"decision": bundle.Decision, "delta_digest": bundle.DeltaDigest, "candidate_digest": bundle.CandidateDigest, "added_cells": bundle.Counts.AddedCells, "retired_cells": bundle.Counts.RetiredCells, "split_cells": bundle.Counts.SplitCells, "rollback_exact_pair": bundle.RollbackDelta.ExactPair},
		"baseline_corpus":          baseline,
		"candidate_corpus":         candidate,
		"causal_decision":          causalDecision,
		"causal_full_oracle_state": causalOracleState,
		"acceptance_predicates":    acceptance,
		"normalization_evidence":   evidence,
		"process_evidence":         process,
		"metrics":                  metrics,
		"unresolved":               []string{"released compiler phase parser accepts exactly three executable activities and cannot compose the requested four-activity split", "candidate semantic IR and backend were therefore not emitted", "whole-language self-improvement remains UNKNOWN", "external utility remains UNKNOWN/NOT_MADE"},
	}
	if err := WriteJSON(reportPath, final); err != nil {
		return err
	}
	return os.WriteFile(dossierPath, []byte(renderDossier(final, baseline, candidate, bundle, causal, evidence, process, metrics)), 0o644)
}

func allGenerated(report CorpusReport) bool {
	if len(report.Cases) != 3 {
		return false
	}
	for _, item := range report.Cases {
		if !item.GeneratedArtifacts {
			return false
		}
	}
	return true
}

func renderDossier(final map[string]any, baseline, candidate CorpusReport, bundle CandidateBundle, causal, evidence, process, metrics map[string]any) string {
	var out strings.Builder
	fmt.Fprintf(&out, "# Gooo evolution trial dossier\n\n")
	fmt.Fprintf(&out, "decision: `%s`\n\n", final["decision"])
	out.WriteString("## Experiment boundary\n\n")
	out.WriteString("One experiment was executed against the immutable `gooo-reflexive-compiler-slice v0.1.1` source, the immutable `gooo-language-delta-forge v0.1.2` contract, and the immutable `gooo-causal-verification-runner v0.1.1` contract. All generated outputs were placed under runner-owned temporary storage. Repository and upstream writes were exactly zero; local test executions were exactly zero.\n\n")
	out.WriteString("## Observed behavior and proposed change\n\n")
	out.WriteString("The released normalization activity `NormalizeSource` parses entity/activity declarations, binds stable IDs, sorts the semantic IR by stable ID, refutes duplicate IDs, and records an UNKNOWN when the required entity is absent. The released compiler produced this evidence from its real `CLOSED` corpus execution.\n\n")
	fmt.Fprintf(&out, "The released delta forge returned `%s` with `added_cells=%d`, `retired_cells=%d`, and `split_cells=%d`. Its exact candidate retires the one coarse normalization cell and adds two cells named `ParseSource` and `ValidateStableIDs`; the inverse rollback is exact=%t. The exact observed resolution pair is before=1 and after=2 phase-localization cells, bound to the phase source digest.\n\n", bundle.Decision, bundle.Counts.AddedCells, bundle.Counts.RetiredCells, bundle.Counts.SplitCells, bundle.RollbackDelta.ExactPair)
	out.WriteString("The generated candidate `.gooo` phase therefore has four executable activities: `ParseSource`, `ValidateStableIDs`, `EmitBackend`, and `VerifyReplay`. This is the requested semantic split expressed through metacode; no candidate Go implementation was hand-authored.\n\n")
	out.WriteString("## Composition result\n\n")
	out.WriteString("The released compiler rejected the candidate phase before normalization because its released phase contract requires exactly three executable activities. Consequently, candidate semantic IR and backend artifacts were not emitted. The compiler behavior did not change; the candidate was not adopted.\n\n")
	fmt.Fprintf(&out, "Causal runner decision: `%v`; full-oracle state: `%v`. It selected the affected candidate-corpus test and reused one exact immutable proof for the stable control, then compared both against the full oracle. The candidate corpus failure remains REFUTED evidence rather than being hidden by selective execution.\n\n", causal["decision"], causalOracleState(causal))
	out.WriteString("## Exact corpus comparison\n\n")
	out.WriteString("| run | CLOSED | UNKNOWN | REFUTED | failed commands | replay digests |\n|---|---:|---:|---:|---:|---|\n")
	fmt.Fprintf(&out, "| baseline | %d | %d | %d | %d | %t |\n", baseline.Closed, baseline.Unknown, baseline.Refuted, baseline.Failed, baseline.ReplayDigestMatch)
	fmt.Fprintf(&out, "| candidate | %d | %d | %d | %d | %t |\n\n", candidate.Closed, candidate.Unknown, candidate.Refuted, candidate.Failed, candidate.ReplayDigestMatch)
	out.WriteString("The released baseline corpus was exactly one `CLOSED`, one `UNKNOWN`, and one `REFUTED`. The candidate phase did not preserve that distribution because all three phase probes were rejected by the three-activity interface.\n\n")
	out.WriteString("## Acceptance predicates\n\n")
	out.WriteString("Acceptance is `REFUTED` under the declared precedence because candidate interface composition failed. Replay equality, corpus preservation, causal closure, rollback, and the exact resolution pair are recorded separately; none is upgraded by inference.\n\n")
	out.WriteString("## Unresolved\n\n")
	out.WriteString("- The released compiler needs a released contract that can execute the two-stage normalization split before this experiment can test changed compiler behavior.\n- Candidate IR/backend replay equality is not observable when the released phase parser rejects the candidate.\n- Whole-language self-improvement remains `UNKNOWN`; this experiment covers one bounded compiler phase only.\n- External utility remains `UNKNOWN`/`NOT_MADE`.\n\n")
	fmt.Fprintf(&out, "Process evidence: bootstrap_direct_main=%v, post_bootstrap_direct_main=%v, repository_writes_by_experiment=%v, upstream_writes_by_experiment=%v.\n", process["bootstrap_direct_main"], process["post_bootstrap_direct_main"], process["repository_writes_by_experiment"], process["upstream_writes_by_experiment"])
	fmt.Fprintf(&out, "Normalization evidence digest: %v.\n", evidence["phase_digest"])
	fmt.Fprintf(&out, "CI metrics receipt: compile_wall_ms=%v build_wall_ms=%v test_wall_ms=%v conformance_wall_ms=%v integration_wall_ms=%v peak_rss_kib=%v generated_files=%v generated_bytes=%v.\n", metrics["compile_wall_ms"], metrics["build_wall_ms"], metrics["test_wall_ms"], metrics["conformance_wall_ms"], metrics["integration_wall_ms"], metrics["peak_rss_kib"], metrics["generated_files"], metrics["generated_bytes"])
	return out.String()
}

func causalOracleState(value map[string]any) string {
	if comparison, ok := value["full_oracle_comparison"].(map[string]any); ok {
		if state, ok := comparison["state"].(string); ok {
			return state
		}
	}
	return "UNKNOWN"
}
