package trial

type DeltaInput struct {
	Schema   string         `json:"schema"`
	Release  DeltaRelease   `json:"release"`
	Failure  DeltaFailure   `json:"failure"`
	Receipts []DeltaReceipt `json:"receipts"`
}

type DeltaRelease struct {
	Schema   string             `json:"schema"`
	Version  string             `json:"version"`
	Identity DeltaReleaseID     `json:"identity"`
	Graph    DeltaSemanticGraph `json:"graph"`
}

type DeltaReleaseID struct {
	Repository   string `json:"repository"`
	Tag          string `json:"tag"`
	Commit       string `json:"commit"`
	SourceDigest string `json:"source_digest"`
	GraphDigest  string `json:"graph_digest"`
	Digest       string `json:"digest"`
}

type DeltaSemanticGraph struct {
	Schema      string           `json:"schema"`
	Version     string           `json:"version"`
	Concepts    []DeltaConcept   `json:"concepts"`
	Predicates  []DeltaPredicate `json:"predicates"`
	Edges       []DeltaGraphEdge `json:"edges"`
	Cells       []DeltaGraphCell `json:"cells"`
	GraphDigest string           `json:"graph_digest"`
}

type DeltaGraphSnapshot struct {
	Digest     string           `json:"digest"`
	Concepts   []DeltaConcept   `json:"concepts"`
	Predicates []DeltaPredicate `json:"predicates"`
	Edges      []DeltaGraphEdge `json:"edges"`
	Cells      []DeltaGraphCell `json:"cells"`
}

type DeltaConcept struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DeltaPredicate struct {
	ID        string       `json:"id"`
	ConceptID string       `json:"concept_id"`
	Name      string       `json:"name"`
	Fields    []DeltaField `json:"fields"`
}

type DeltaField struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type DeltaGraphEdge struct {
	ID          string `json:"id"`
	From        string `json:"from"`
	To          string `json:"to"`
	Kind        string `json:"kind"`
	PredicateID string `json:"predicate_id"`
}

type DeltaGraphCell struct {
	ID          string `json:"id"`
	ConceptID   string `json:"concept_id"`
	PredicateID string `json:"predicate_id"`
	FieldID     string `json:"field_id"`
	Relation    string `json:"relation"`
	Constraint  string `json:"constraint"`
}

type DeltaFailure struct {
	Schema         string            `json:"schema"`
	ID             string            `json:"id"`
	State          string            `json:"state"`
	RelatedTo      string            `json:"related_to,omitempty"`
	SourceDigest   string            `json:"source_digest"`
	GraphDigest    string            `json:"graph_digest"`
	Stage          string            `json:"stage,omitempty"`
	Step           string            `json:"step,omitempty"`
	Reason         string            `json:"reason,omitempty"`
	UnknownClass   string            `json:"unknown_class,omitempty"`
	NextOperation  string            `json:"next_operation,omitempty"`
	BlockedBy      []string          `json:"blocked_by,omitempty"`
	Target         *DeltaTarget      `json:"target,omitempty"`
	DirectCause    *DeltaDirectCause `json:"direct_cause,omitempty"`
	Counterexample string            `json:"counterexample,omitempty"`
	Digest         string            `json:"digest,omitempty"`
}

type DeltaReceipt struct {
	Schema         string   `json:"schema"`
	ID             string   `json:"id"`
	State          string   `json:"state"`
	RelatedTo      string   `json:"related_to,omitempty"`
	SourceDigest   string   `json:"source_digest"`
	GraphDigest    string   `json:"graph_digest"`
	Stage          string   `json:"stage,omitempty"`
	Step           string   `json:"step,omitempty"`
	Reason         string   `json:"reason,omitempty"`
	UnknownClass   string   `json:"unknown_class,omitempty"`
	NextOperation  string   `json:"next_operation,omitempty"`
	BlockedBy      []string `json:"blocked_by,omitempty"`
	Counterexample string   `json:"counterexample,omitempty"`
	Digest         string   `json:"digest,omitempty"`
}

type DeltaTarget struct {
	ConceptID   string `json:"concept_id"`
	PredicateID string `json:"predicate_id"`
	FieldID     string `json:"field_id"`
}

type DeltaDirectCause struct {
	Target               DeltaTarget             `json:"target"`
	CausalFrontier       []DeltaFrontierEdge     `json:"causal_frontier"`
	BeforeGraphDigest    string                  `json:"before_graph_digest"`
	AfterGraph           DeltaGraphSnapshot      `json:"after_graph"`
	AffectedPredicateIDs []string                `json:"affected_predicate_ids"`
	AddedCells           []DeltaGraphCell        `json:"added_cells"`
	RetiredCells         []DeltaGraphCell        `json:"retired_cells"`
	SplitCells           []DeltaSplitCell        `json:"split_cells"`
	Rollback             DeltaRollback           `json:"rollback"`
	PositiveCases        []DeltaConformanceCase  `json:"positive_cases"`
	NegativeCases        []DeltaConformanceCase  `json:"negative_cases"`
	IndependentConsumer  DeltaIndependentReceipt `json:"independent_consumer"`
}

type DeltaFrontierEdge struct {
	EdgeID      string `json:"edge_id"`
	From        string `json:"from"`
	To          string `json:"to"`
	Kind        string `json:"kind"`
	PredicateID string `json:"predicate_id"`
}

type DeltaSplitCell struct {
	SplitID       string   `json:"split_id"`
	RetiredCellID string   `json:"retired_cell_id"`
	AddedCellIDs  []string `json:"added_cell_ids"`
}

type DeltaRollback struct {
	RemoveAddedCellIDs  []string         `json:"remove_added_cell_ids"`
	RestoreRetiredCells []DeltaGraphCell `json:"restore_retired_cells"`
	Unsplit             []DeltaSplitCell `json:"unsplit"`
	ExactPair           bool             `json:"exact_pair"`
}

type DeltaConformanceCase struct {
	ID                     string      `json:"id"`
	Polarity               string      `json:"polarity"`
	BeforeGraphDigest      string      `json:"before_graph_digest"`
	AfterGraphDigest       string      `json:"after_graph_digest"`
	Target                 DeltaTarget `json:"target"`
	ExpectedPredicateState string      `json:"expected_predicate_state"`
	FixtureDigest          string      `json:"fixture_digest"`
}

type DeltaIndependentReceipt struct {
	Schema              string   `json:"schema"`
	ReceiptID           string   `json:"receipt_id"`
	ProducerID          string   `json:"producer_id"`
	ConsumerID          string   `json:"consumer_id"`
	State               string   `json:"state"`
	SourceDigest        string   `json:"source_digest,omitempty"`
	BaselineGraphDigest string   `json:"baseline_graph_digest,omitempty"`
	DeltaDigest         string   `json:"delta_digest,omitempty"`
	ReceiptDigest       string   `json:"receipt_digest,omitempty"`
	Reason              string   `json:"reason,omitempty"`
	BlockedBy           []string `json:"blocked_by,omitempty"`
}

type CandidateBundle struct {
	Schema               string                  `json:"schema"`
	Version              string                  `json:"version"`
	Decision             string                  `json:"decision"`
	SourceDigest         string                  `json:"source_digest"`
	BaselineGraphDigest  string                  `json:"baseline_graph_digest"`
	DeltaDigest          string                  `json:"delta_digest"`
	Target               DeltaTargetResolution   `json:"target"`
	CausalFrontier       []DeltaFrontierEdge     `json:"causal_frontier"`
	SemanticGraphDelta   DeltaSemanticGraphDelta `json:"semantic_graph_delta"`
	RollbackDelta        DeltaRollback           `json:"rollback_delta"`
	AffectedPredicateIDs []string                `json:"affected_predicate_ids"`
	Counts               DeltaCounts             `json:"counts"`
	TestManifest         DeltaTestManifest       `json:"test_manifest"`
	Claim                DeltaClaim              `json:"claim"`
	Improvement          DeltaClaim              `json:"improvement"`
	CandidateDigest      string                  `json:"candidate_digest"`
}

type DeltaTargetResolution struct {
	ResolutionLevel string `json:"resolution_level"`
	ConceptID       string `json:"concept_id"`
	PredicateID     string `json:"predicate_id"`
	FieldID         string `json:"field_id"`
}

type DeltaSemanticGraphDelta struct {
	Before       DeltaGraphSnapshot `json:"before"`
	After        DeltaGraphSnapshot `json:"after"`
	AddedCells   []DeltaGraphCell   `json:"added_cells"`
	RetiredCells []DeltaGraphCell   `json:"retired_cells"`
	SplitCells   []DeltaSplitCell   `json:"split_cells"`
	ExactPair    bool               `json:"exact_pair"`
}

type DeltaCounts struct {
	AddedCells   int `json:"added_cells"`
	RetiredCells int `json:"retired_cells"`
	SplitCells   int `json:"split_cells"`
}

type DeltaTestManifest struct {
	ExactBeforeAfterPair bool `json:"exact_before_after_pair"`
}

type DeltaClaim struct {
	State         string   `json:"state"`
	Stage         string   `json:"stage"`
	Step          string   `json:"step"`
	Reason        string   `json:"reason"`
	UnknownClass  string   `json:"unknown_class"`
	NextOperation string   `json:"next_operation"`
	BlockedBy     []string `json:"blocked_by"`
}

type CorpusCase struct {
	ID       string `json:"id"`
	Source   string `json:"source"`
	Expected string `json:"expected"`
}

type CorpusContract struct {
	Cases []CorpusCase `json:"cases"`
}

type CorpusCaseResult struct {
	ID                 string `json:"id"`
	Source             string `json:"source"`
	Expected           string `json:"expected"`
	Observed           string `json:"observed"`
	TerminalResult     string `json:"terminal_result"`
	WallMS             int64  `json:"wall_ms"`
	PeakRSSKiB         int64  `json:"peak_rss_kib"`
	Error              string `json:"error,omitempty"`
	ReplayDigestMatch  bool   `json:"replay_digest_match"`
	GeneratedArtifacts bool   `json:"generated_artifacts"`
}

type CorpusReport struct {
	Schema             string             `json:"schema"`
	Role               string             `json:"role"`
	PhasePath          string             `json:"phase_path"`
	PhaseDigest        string             `json:"phase_digest"`
	CorpusContractPath string             `json:"corpus_contract_path"`
	Cases              []CorpusCaseResult `json:"cases"`
	Total              int                `json:"total"`
	Executed           int                `json:"executed"`
	Reused             int                `json:"reused"`
	Failed             int                `json:"failed"`
	Unknown            int                `json:"unknown"`
	Closed             int                `json:"closed"`
	Refuted            int                `json:"refuted"`
	WallMS             int64              `json:"wall_ms"`
	PeakRSSKiB         int64              `json:"peak_rss_kib"`
	ReplayDigestMatch  bool               `json:"replay_digest_match"`
	InterfaceDecision  string             `json:"interface_decision"`
}

type CausalCase struct {
	Schema          string                 `json:"schema"`
	CaseID          string                 `json:"case_id"`
	Description     string                 `json:"description"`
	ChangeClaim     CausalChangeClaim      `json:"change_claim"`
	SemanticGraph   CausalSemanticGraph    `json:"semantic_graph"`
	Tests           []CausalTestSpec       `json:"tests"`
	ReusableProofs  []CausalProof          `json:"reusable_proofs"`
	FullOracle      CausalFullOracle       `json:"full_oracle"`
	PerformancePair CausalPerformancePair  `json:"performance_pair"`
	Counterexamples []CausalCounterexample `json:"counterexamples"`
	CacheHit        bool                   `json:"cache_hit"`
	Expected        CausalExpected         `json:"expected"`
}

type CausalChangeClaim struct {
	ClaimID           string   `json:"claim_id"`
	ChangedPredicates []string `json:"changed_predicates"`
	SourceDigest      string   `json:"source_digest"`
	SemanticDigest    string   `json:"semantic_digest"`
}

type CausalSemanticGraph struct {
	Schema         string            `json:"schema"`
	GraphID        string            `json:"graph_id"`
	SourceDigest   string            `json:"source_digest"`
	SemanticDigest string            `json:"semantic_digest"`
	GraphDigest    string            `json:"graph_digest"`
	Predicates     []string          `json:"predicates"`
	Edges          []CausalGraphEdge `json:"edges"`
}

type CausalGraphEdge struct {
	EdgeID string `json:"edge_id"`
	From   string `json:"from"`
	To     string `json:"to"`
	Kind   string `json:"kind"`
	Status string `json:"status"`
}

type CausalTestSpec struct {
	TestID      string             `json:"test_id"`
	TargetNode  string             `json:"target_node"`
	Command     []string           `json:"command"`
	Observation *CausalObservation `json:"observation,omitempty"`
}

type CausalObservation struct {
	Status       string `json:"status"`
	WallMS       int64  `json:"wall_ms"`
	PeakRSSKiB   int64  `json:"peak_rss_kib"`
	ResultDigest string `json:"result_digest"`
}

type CausalProof struct {
	ProofID             string                `json:"proof_id"`
	TestID              string                `json:"test_id"`
	SourceDigest        string                `json:"source_digest"`
	SourceTreeDigest    string                `json:"source_tree_digest"`
	SemanticDigest      string                `json:"semantic_digest"`
	GraphDigest         string                `json:"graph_digest"`
	ToolchainDigest     string                `json:"toolchain_digest"`
	ScenarioDigest      string                `json:"scenario_digest"`
	TestInventoryDigest string                `json:"test_inventory_digest"`
	CommandDigest       string                `json:"command_digest"`
	TerminalResult      string                `json:"terminal_result"`
	ResultDigest        string                `json:"result_digest"`
	Release             CausalReleaseEvidence `json:"release"`
}

type CausalReleaseEvidence struct {
	Provider          string `json:"provider"`
	Repository        string `json:"repository"`
	ReleaseID         int64  `json:"release_id"`
	Tag               string `json:"tag"`
	SelfAsserted      bool   `json:"self_asserted_immutable"`
	PlatformImmutable bool   `json:"platform_immutable"`
	AssetDigest       string `json:"asset_digest"`
	AssetName         string `json:"asset_name"`
}

type CausalFullOracle struct {
	OracleID    string               `json:"oracle_id"`
	Independent bool                 `json:"independent"`
	Digest      string               `json:"digest"`
	Results     []CausalOracleResult `json:"results"`
}

type CausalOracleResult struct {
	TestID       string `json:"test_id"`
	Status       string `json:"status"`
	ResultDigest string `json:"result_digest"`
}

type CausalPerformancePair struct {
	Before CausalTimingSnapshot `json:"before"`
	After  CausalTimingSnapshot `json:"after"`
}

type CausalTimingSnapshot struct {
	SourceDigest    string `json:"source_digest"`
	ToolchainDigest string `json:"toolchain_digest"`
	ScenarioDigest  string `json:"scenario_digest"`
	WallMS          int64  `json:"wall_ms"`
	PeakRSSKiB      int64  `json:"peak_rss_kib"`
}

type CausalCounterexample struct {
	CounterexampleID     string `json:"counterexample_id"`
	TestID               string `json:"test_id"`
	ExpectedInvalidation bool   `json:"expected_invalidation"`
	ObservedInvalidation bool   `json:"observed_invalidation"`
}

type CausalExpected struct {
	Decision             string `json:"decision"`
	SelectionMode        string `json:"selection_mode"`
	TotalTests           int    `json:"total_tests"`
	Selected             int    `json:"selected"`
	Executed             int    `json:"executed"`
	Reused               int    `json:"reused"`
	FullOracleExecuted   int    `json:"full_oracle_executed"`
	Failures             int    `json:"failures"`
	Unknowns             int    `json:"unknowns"`
	AvoidedExecutions    int    `json:"avoided_executions"`
	InvalidatedEdgeCount int    `json:"invalidated_edge_count"`
}

type AuthorityReport struct {
	Schema              string         `json:"schema"`
	Decision            string         `json:"decision"`
	Activities          int            `json:"activities"`
	Cells               int            `json:"cells"`
	ProofTotals         map[string]int `json:"proof_totals"`
	IndicatorTotals     map[string]int `json:"indicator_totals"`
	Precedence          []string       `json:"precedence"`
	UnknownFields       []string       `json:"unknown_fields"`
	RepositoryWrites    int            `json:"repository_writes"`
	UpstreamWrites      int            `json:"upstream_writes"`
	LocalTestExecutions int            `json:"local_test_executions"`
	Errors              []string       `json:"errors"`
}
