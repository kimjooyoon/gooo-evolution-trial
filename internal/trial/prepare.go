package trial

import (
	"fmt"
)

const (
	deltaInputSchema       = "gooo/language-delta-forge/input/v1"
	deltaReleaseSchema     = "gooo/language-delta-forge/immutable-release/v1"
	deltaReceiptSchema     = "gooo/language-delta-forge/observation-receipt/v1"
	deltaIndependentSchema = "gooo/language-delta-forge/independent-consumer-receipt/v1"
)

func PrepareDeltaInput(phasePath, compilerCommit, baselineReceiptPath, outputDir string) error {
	if compilerCommit == "" {
		return fmt.Errorf("compiler commit is required")
	}
	if err := PrepareEmpty(outputDir); err != nil {
		return err
	}
	sourceDigest, sourceData, err := DigestFile(phasePath)
	if err != nil {
		return err
	}
	graph := reflexiveGraph()
	graph.GraphDigest, err = digestDeltaSnapshot(snapshotForGraph(graph))
	if err != nil {
		return err
	}
	after := splitGraph(graph)
	after.GraphDigest, err = digestDeltaSnapshot(snapshotForGraph(after))
	if err != nil {
		return err
	}
	identity := DeltaReleaseID{
		Repository: "kimjooyoon/gooo-reflexive-compiler-slice",
		Tag:        "v0.1.1", Commit: compilerCommit, SourceDigest: sourceDigest, GraphDigest: graph.GraphDigest,
	}
	identity.Digest, err = digestDeltaIdentity(identity)
	if err != nil {
		return err
	}
	target := DeltaTarget{ConceptID: "concept.reflexive-compiler", PredicateID: "predicate.normalize-source", FieldID: "field.normalize-source.phase"}
	added := []DeltaGraphCell{after.Cells[0], after.Cells[1]}
	retired := []DeltaGraphCell{graph.Cells[0]}
	split := []DeltaSplitCell{{SplitID: "split.normalize-source.coarse", RetiredCellID: retired[0].ID, AddedCellIDs: []string{added[0].ID, added[1].ID}}}
	rollback := DeltaRollback{RemoveAddedCellIDs: []string{added[0].ID, added[1].ID}, RestoreRetiredCells: retired, Unsplit: split, ExactPair: true}
	positive := []DeltaConformanceCase{{
		ID: "positive-normalize-source-split", Polarity: "POSITIVE", BeforeGraphDigest: graph.GraphDigest, AfterGraphDigest: after.GraphDigest,
		Target: target, ExpectedPredicateState: "CLOSED", FixtureDigest: "sha256:1111111111111111111111111111111111111111111111111111111111111111",
	}}
	negative := []DeltaConformanceCase{{
		ID: "negative-normalize-source-split", Polarity: "NEGATIVE", BeforeGraphDigest: graph.GraphDigest, AfterGraphDigest: after.GraphDigest,
		Target: target, ExpectedPredicateState: "REFUTED", FixtureDigest: "sha256:2222222222222222222222222222222222222222222222222222222222222222",
	}}
	direct := DeltaDirectCause{
		Target:            target,
		CausalFrontier:    []DeltaFrontierEdge{{EdgeID: "edge.normalize-source", From: "concept.reflexive-compiler", To: "predicate.normalize-source", Kind: "binds", PredicateID: "predicate.normalize-source"}},
		BeforeGraphDigest: graph.GraphDigest, AfterGraph: snapshotForGraph(after), AffectedPredicateIDs: []string{"predicate.normalize-source"},
		AddedCells: added, RetiredCells: retired, SplitCells: split, Rollback: rollback, PositiveCases: positive, NegativeCases: negative,
		IndependentConsumer: DeltaIndependentReceipt{Schema: deltaIndependentSchema, ReceiptID: "independent-normalize-source", ProducerID: "reflexive-observer", ConsumerID: "evolution-trial-independent-consumer", State: "CLOSED"},
	}
	failure := DeltaFailure{Schema: deltaReceiptSchema, ID: "failure-reflexive-normalize-source", State: "FAILURE", SourceDigest: sourceDigest, GraphDigest: graph.GraphDigest, Target: &target, DirectCause: &direct}
	input := DeltaInput{Schema: deltaInputSchema, Release: DeltaRelease{Schema: deltaReleaseSchema, Version: "v1", Identity: identity, Graph: graph}, Failure: failure, Receipts: []DeltaReceipt{}}
	if err := WriteJSON(outputDir+"/input.json", input); err != nil {
		return err
	}
	if err := WriteJSON(outputDir+"/release.json", input.Release); err != nil {
		return err
	}
	evidence := map[string]any{
		"schema":     "gooo/evolution-trial/normalization-evidence/v1",
		"phase_path": phasePath, "phase_digest": sourceDigest, "phase_bytes": len(sourceData),
		"observed_activity": "NormalizeSource", "observed_program": "reflexive.normalize:v1",
		"resolution_pair":  map[string]any{"before": 1, "after": 2, "unit": "phase-localization-cells", "matched_source_digest": sourceDigest, "matched_toolchain": "github-actions-go-1.27"},
		"baseline_receipt": baselineReceiptPath,
	}
	if baselineReceiptPath != "" {
		receiptDigest, _, digestErr := DigestFile(baselineReceiptPath)
		if digestErr != nil {
			return digestErr
		}
		evidence["baseline_receipt_digest"] = receiptDigest
	}
	return WriteJSON(outputDir+"/normalization-evidence.json", evidence)
}

func reflexiveGraph() DeltaSemanticGraph {
	return DeltaSemanticGraph{
		Schema: "gooo/semantic-graph/v1", Version: "v1",
		Concepts: []DeltaConcept{{ID: "concept.reflexive-compiler", Name: "ReflexiveCompiler"}},
		Predicates: []DeltaPredicate{
			{ID: "predicate.normalize-source", ConceptID: "concept.reflexive-compiler", Name: "NormalizeSource", Fields: []DeltaField{{ID: "field.normalize-source.phase", Name: "phase", Type: "string"}}},
			{ID: "predicate.emit-backend", ConceptID: "concept.reflexive-compiler", Name: "EmitBackend", Fields: []DeltaField{{ID: "field.emit-backend.artifact", Name: "artifact", Type: "string"}}},
			{ID: "predicate.verify-replay", ConceptID: "concept.reflexive-compiler", Name: "VerifyReplay", Fields: []DeltaField{{ID: "field.verify-replay.rollback", Name: "rollback", Type: "string"}}},
		},
		Edges: []DeltaGraphEdge{
			{ID: "edge.normalize-source", From: "concept.reflexive-compiler", To: "predicate.normalize-source", Kind: "binds", PredicateID: "predicate.normalize-source"},
			{ID: "edge.emit-backend", From: "concept.reflexive-compiler", To: "predicate.emit-backend", Kind: "binds", PredicateID: "predicate.emit-backend"},
			{ID: "edge.verify-replay", From: "concept.reflexive-compiler", To: "predicate.verify-replay", Kind: "binds", PredicateID: "predicate.verify-replay"},
		},
		Cells: []DeltaGraphCell{
			{ID: "cell.normalize-source.coarse", ConceptID: "concept.reflexive-compiler", PredicateID: "predicate.normalize-source", FieldID: "field.normalize-source.phase", Relation: "requires", Constraint: "one-coarse-phase"},
			{ID: "cell.emit-backend.derived", ConceptID: "concept.reflexive-compiler", PredicateID: "predicate.emit-backend", FieldID: "field.emit-backend.artifact", Relation: "requires", Constraint: "backend-only"},
			{ID: "cell.verify-replay.rollback", ConceptID: "concept.reflexive-compiler", PredicateID: "predicate.verify-replay", FieldID: "field.verify-replay.rollback", Relation: "requires", Constraint: "retain-baseline"},
		},
	}
}

func splitGraph(before DeltaSemanticGraph) DeltaSemanticGraph {
	after := before
	after.Edges = append([]DeltaGraphEdge(nil), before.Edges...)
	after.Cells = []DeltaGraphCell{
		{ID: "cell.normalize-source.parse-source", ConceptID: "concept.reflexive-compiler", PredicateID: "predicate.normalize-source", FieldID: "field.normalize-source.phase", Relation: "requires", Constraint: "ParseSource"},
		{ID: "cell.normalize-source.validate-stable-ids", ConceptID: "concept.reflexive-compiler", PredicateID: "predicate.normalize-source", FieldID: "field.normalize-source.phase", Relation: "requires", Constraint: "ValidateStableIDs"},
		before.Cells[1], before.Cells[2],
	}
	after.Edges = append(after.Edges,
		DeltaGraphEdge{ID: "edge.normalize-source.parse-source", From: "predicate.normalize-source", To: "cell.normalize-source.parse-source", Kind: "refines", PredicateID: "predicate.normalize-source"},
		DeltaGraphEdge{ID: "edge.normalize-source.validate-stable-ids", From: "predicate.normalize-source", To: "cell.normalize-source.validate-stable-ids", Kind: "refines", PredicateID: "predicate.normalize-source"},
	)
	return after
}

func snapshotForGraph(graph DeltaSemanticGraph) DeltaGraphSnapshot {
	return DeltaGraphSnapshot{Digest: graph.GraphDigest, Concepts: graph.Concepts, Predicates: graph.Predicates, Edges: graph.Edges, Cells: graph.Cells}
}
