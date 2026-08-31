package trial

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type authorityContract struct {
	Schema           string          `json:"schema"`
	DenominatorID    string          `json:"denominator_id"`
	FixedDenominator int             `json:"fixed_denominator"`
	Precedence       []string        `json:"precedence"`
	UnknownFields    []string        `json:"unknown_fields"`
	ProofTotals      map[string]int  `json:"proof_totals"`
	IndicatorTotals  map[string]int  `json:"indicator_totals"`
	Cells            []authorityCell `json:"cells"`
}

type authorityCell struct {
	Ordinal        int    `json:"ordinal"`
	ID             string `json:"id"`
	Activity       string `json:"activity"`
	Stage          string `json:"stage"`
	Step           string `json:"step"`
	ProofChoice    string `json:"proof_choice"`
	IndicatorClass string `json:"indicator_class"`
}

func VerifyAuthority(metaPath, contractPath string) (AuthorityReport, error) {
	report := AuthorityReport{
		Schema:        "gooo/evolution-trial/authority-report/v1",
		Precedence:    []string{"REFUTED", "UNKNOWN", "CLOSED"},
		UnknownFields: []string{"stage", "step", "reason", "unknown_class", "next_operation", "blocked_by"},
		ProofTotals:   map[string]int{}, IndicatorTotals: map[string]int{}, Errors: []string{},
	}
	meta, err := os.ReadFile(metaPath)
	if err != nil {
		return report, err
	}
	var contract authorityContract
	if err := ReadJSON(contractPath, &contract); err != nil {
		return report, err
	}
	text := string(meta)
	markers := []string{
		"gooo evolution_trial v1",
		"precedence REFUTED>UNKNOWN>CLOSED",
		"unknown_fields stage,step,reason,unknown_class,next_operation,blocked_by",
		"authority repository_writes=0 upstream_writes=0 local_test_executions=0 merge_authority=0",
		"metrics exact-integers-only",
		"adoption candidate-only automatic-merge=false protected-core-mutation=false",
		"rollback baseline-retained candidate-separate inverse-delta-required",
		"process_field transition \"PR_MERGE_ONLY\"",
	}
	for _, marker := range markers {
		if !strings.Contains(text, marker) {
			report.Errors = append(report.Errors, "missing meta authority: "+marker)
		}
	}
	if contract.Schema != "gooo/evolution-trial/denominator/v1" || contract.DenominatorID != "evolution-trial-v1" || contract.FixedDenominator != 18 {
		report.Errors = append(report.Errors, "fixed denominator is not 18")
	}
	if !equalStrings(contract.Precedence, report.Precedence) {
		report.Errors = append(report.Errors, "precedence contract mismatch")
	}
	if !equalStrings(contract.UnknownFields, report.UnknownFields) {
		report.Errors = append(report.Errors, "UNKNOWN tuple contract mismatch")
	}
	if len(contract.Cells) != 18 {
		report.Errors = append(report.Errors, "denominator cell count is not 18")
	}
	activityNames := make([]string, 0)
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "activity ") {
			continue
		}
		report.Activities++
		body := strings.TrimPrefix(line, "activity ")
		nameEnd := strings.IndexByte(body, '(')
		if nameEnd > 0 {
			activityNames = append(activityNames, body[:nameEnd])
		}
		quoted := strings.Split(body, "\"")
		if len(quoted) >= 2 {
			for _, part := range strings.Split(quoted[1], ";") {
				key, value, ok := strings.Cut(part, "=")
				if !ok {
					continue
				}
				switch key {
				case "proof":
					report.ProofTotals[value]++
				case "indicator":
					report.IndicatorTotals[value]++
				}
			}
		}
	}
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "cell ") {
			report.Cells++
		}
	}
	if report.Activities != 18 || report.Cells != 18 {
		report.Errors = append(report.Errors, fmt.Sprintf("meta activity/cell count is %d/%d, expected 18/18", report.Activities, report.Cells))
	}
	for _, proof := range []string{"FOUNDATION", "COHERENCE", "REGRESSION"} {
		if report.ProofTotals[proof] != 6 || contract.ProofTotals[proof] != 6 {
			report.Errors = append(report.Errors, "proof denominator mismatch: "+proof)
		}
	}
	for _, indicator := range []string{"DRIVER", "OUTCOME", "GUARDRAIL"} {
		if report.IndicatorTotals[indicator] != 6 || contract.IndicatorTotals[indicator] != 6 {
			report.Errors = append(report.Errors, "indicator denominator mismatch: "+indicator)
		}
	}
	if len(activityNames) != len(contract.Cells) {
		report.Errors = append(report.Errors, "activity and denominator topology mismatch")
	} else {
		for index, cell := range contract.Cells {
			if cell.Ordinal != index+1 || cell.Activity != activityNames[index] {
				report.Errors = append(report.Errors, "cell activity mismatch at ordinal "+strconv.Itoa(index+1))
			}
		}
	}
	if !strings.Contains(text, "process_field repository \"kimjooyoon/gooo-evolution-trial\"") || !strings.Contains(text, "process_evidence_schema \"gooo/evolution-trial/process-evidence/v1\"") {
		report.Errors = append(report.Errors, "process authority evidence binding is missing")
	}
	report.RepositoryWrites = 0
	report.UpstreamWrites = 0
	report.LocalTestExecutions = 0
	if len(report.Errors) == 0 {
		report.Decision = "CLOSED"
	} else {
		report.Decision = "REFUTED"
	}
	return report, nil
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
