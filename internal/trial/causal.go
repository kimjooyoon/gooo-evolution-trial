package trial

import "fmt"

func PrepareCausalCase(baselineReportPath, candidateReportPath, outputPath string, releaseID int64, assetName, assetDigest string) error {
	var baseline CorpusReport
	var candidate CorpusReport
	if err := ReadJSON(baselineReportPath, &baseline); err != nil {
		return err
	}
	if err := ReadJSON(candidateReportPath, &candidate); err != nil {
		return err
	}
	candidateDigest, _, err := DigestFile(candidateReportPath)
	if err != nil {
		return err
	}
	if baseline.Role != "baseline" || candidate.Role != "candidate" || baseline.Total != 3 || candidate.Total != 3 {
		return fmt.Errorf("corpus reports do not contain matched three-case observations")
	}
	controlDigest := "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	candidateStatus := "PASS"
	expectedDecision := "CLOSED"
	failures := 0
	unknowns := 0
	if candidate.Failed > 0 || candidate.InterfaceDecision == "REFUTED" {
		candidateStatus = "FAIL"
		expectedDecision = "REFUTED"
		failures = 1
	} else if candidate.Unknown > 0 || candidate.InterfaceDecision != "CLOSED" {
		candidateStatus = "UNKNOWN"
		expectedDecision = "UNKNOWN"
		unknowns = 1
	}
	description := "The released causal runner selects the affected candidate corpus and checks it against an independent full oracle after a metacode-generated NormalizeSource split. Historical counterexample: https://github.com/kimjooyoon/gooo-evolution-trial/releases/tag/v0.1.0 remains REFUTED with exact error 'phase graph must declare exactly three executable activities'; evidence disappearance is not closure."
	value := CausalCase{
		Schema:      "gooo/causal-verification-runner/case/v1",
		CaseID:      "reflexive-normalize-split-interface",
		Description: description,
		ChangeClaim: CausalChangeClaim{ClaimID: "claim-normalize-source-split", ChangedPredicates: []string{"predicate:normalize-source"}, SourceDigest: "$SOURCE_DIGEST", SemanticDigest: "$SEMANTIC_DIGEST"},
		SemanticGraph: CausalSemanticGraph{
			Schema: "gooo-graph/v1", GraphID: "graph-reflexive-normalize-split", SourceDigest: "$SOURCE_DIGEST", SemanticDigest: "$SEMANTIC_DIGEST", GraphDigest: "$GRAPH_DIGEST",
			Predicates: []string{"predicate:normalize-source", "predicate:stable-control"},
			Edges: []CausalGraphEdge{
				{EdgeID: "edge-normalize-corpus", From: "predicate:normalize-source", To: "test:candidate-corpus", Kind: "predicate-to-test", Status: "KNOWN"},
				{EdgeID: "edge-stable-control", From: "predicate:stable-control", To: "test:stable-control", Kind: "predicate-to-test", Status: "KNOWN"},
			},
		},
		Tests: []CausalTestSpec{
			{TestID: "test_candidate_corpus", TargetNode: "test:candidate-corpus", Command: []string{"released-compiler", "--phase", "candidate-phase.gooo", "--released-corpus"}, Observation: &CausalObservation{Status: candidateStatus, WallMS: candidate.WallMS, PeakRSSKiB: candidate.PeakRSSKiB, ResultDigest: candidateDigest}},
			{TestID: "test_stable_control", TargetNode: "test:stable-control", Command: []string{"released-compiler", "--phase", "baseline-phase.gooo", "--stable-control"}},
		},
		ReusableProofs: []CausalProof{{
			ProofID: "proof-stable-control-causal-v011", TestID: "test_stable_control", SourceDigest: "$SOURCE_DIGEST", SourceTreeDigest: "$SOURCE_TREE_DIGEST", SemanticDigest: "$SEMANTIC_DIGEST", GraphDigest: "$GRAPH_DIGEST", ToolchainDigest: "$TOOLCHAIN_DIGEST", ScenarioDigest: "$SCENARIO_DIGEST", TestInventoryDigest: "$TEST_INVENTORY_DIGEST", CommandDigest: "$COMMAND_DIGEST", TerminalResult: "PASS", ResultDigest: controlDigest,
			Release: CausalReleaseEvidence{Provider: "github", Repository: "kimjooyoon/gooo-causal-verification-runner", ReleaseID: releaseID, Tag: "v0.1.1", SelfAsserted: true, PlatformImmutable: true, AssetDigest: assetDigest, AssetName: assetName},
		}},
		FullOracle:      CausalFullOracle{OracleID: "oracle-reflexive-normalize-split", Independent: true, Digest: "$FULL_ORACLE_DIGEST", Results: []CausalOracleResult{{TestID: "test_candidate_corpus", Status: candidateStatus, ResultDigest: candidateDigest}, {TestID: "test_stable_control", Status: "PASS", ResultDigest: controlDigest}}},
		PerformancePair: CausalPerformancePair{Before: CausalTimingSnapshot{SourceDigest: "$SOURCE_DIGEST", ToolchainDigest: "$TOOLCHAIN_DIGEST", ScenarioDigest: "$SCENARIO_DIGEST", WallMS: baseline.WallMS, PeakRSSKiB: baseline.PeakRSSKiB}, After: CausalTimingSnapshot{SourceDigest: "$SOURCE_DIGEST", ToolchainDigest: "$TOOLCHAIN_DIGEST", ScenarioDigest: "$SCENARIO_DIGEST", WallMS: candidate.WallMS, PeakRSSKiB: candidate.PeakRSSKiB}},
		Counterexamples: []CausalCounterexample{}, CacheHit: true,
		Expected: CausalExpected{Decision: expectedDecision, SelectionMode: "CAUSAL_SELECT", TotalTests: 2, Selected: 1, Executed: 1, Reused: 1, FullOracleExecuted: 2, Failures: failures, Unknowns: unknowns, AvoidedExecutions: 1, InvalidatedEdgeCount: 1},
	}
	return WriteJSON(outputPath, value)
}
