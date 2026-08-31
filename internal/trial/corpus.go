package trial

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type CorpusOptions struct {
	Compiler   string
	Verifier   string
	Phase      string
	Root       string
	Role       string
	OutputDir  string
	ReportPath string
}

type processObservation struct {
	Status     int
	WallMS     int64
	PeakRSSKiB int64
	Error      string
}

func RunCorpus(options CorpusOptions) (CorpusReport, error) {
	report := CorpusReport{Schema: "gooo/evolution-trial/corpus-report/v1", Role: options.Role, PhasePath: options.Phase, CorpusContractPath: filepath.Join(options.Root, "contracts", "denominator-v1.json"), Cases: []CorpusCaseResult{}, Reused: 0}
	if options.Role != "baseline" && options.Role != "candidate" {
		return report, fmt.Errorf("role must be baseline or candidate")
	}
	if err := PrepareEmpty(options.OutputDir); err != nil {
		return report, err
	}
	phaseDigest, _, err := DigestFile(options.Phase)
	if err != nil {
		return report, err
	}
	report.PhaseDigest = phaseDigest
	var contract CorpusContract
	if err := ReadJSON(report.CorpusContractPath, &contract); err != nil {
		return report, err
	}
	if len(contract.Cases) != 3 {
		return report, fmt.Errorf("released corpus must contain exactly three cases")
	}
	start := time.Now()
	for _, item := range contract.Cases {
		caseResult := CorpusCaseResult{ID: item.ID, Source: item.Source, Expected: item.Expected}
		caseDir := filepath.Join(options.OutputDir, item.ID)
		if err := os.MkdirAll(caseDir, 0o755); err != nil {
			return report, err
		}
		sourcePath := filepath.Join(options.Root, item.Source)
		if options.Role == "baseline" {
			caseResult = runBaselineCase(options, item, sourcePath, caseDir, caseResult)
		} else {
			caseResult = runCandidateCase(options, item, sourcePath, caseDir, caseResult)
		}
		report.Cases = append(report.Cases, caseResult)
		report.Total++
		report.Executed++
		report.WallMS += caseResult.WallMS
		if caseResult.PeakRSSKiB > report.PeakRSSKiB {
			report.PeakRSSKiB = caseResult.PeakRSSKiB
		}
		switch caseResult.Observed {
		case "CLOSED":
			report.Closed++
		case "UNKNOWN":
			report.Unknown++
		case "REFUTED":
			report.Refuted++
		}
		if caseResult.TerminalResult != "PASS" || caseResult.Observed != caseResult.Expected {
			report.Failed++
		}
		if !caseResult.ReplayDigestMatch {
			report.ReplayDigestMatch = false
		}
	}
	if len(report.Cases) > 0 {
		report.ReplayDigestMatch = true
		for _, caseResult := range report.Cases {
			if !caseResult.ReplayDigestMatch {
				report.ReplayDigestMatch = false
			}
		}
	}
	if options.Role == "baseline" && report.Failed == 0 && report.Closed == 1 && report.Unknown == 1 && report.Refuted == 1 && report.ReplayDigestMatch {
		report.InterfaceDecision = "CLOSED"
	} else if options.Role == "candidate" && report.Failed == report.Total && report.Total == 3 {
		report.InterfaceDecision = "REFUTED"
	} else {
		report.InterfaceDecision = "UNKNOWN"
	}
	report.WallMS = time.Since(start).Milliseconds()
	if report.WallMS < 1 {
		report.WallMS = 1
	}
	if options.ReportPath != "" {
		if err := WriteJSON(options.ReportPath, report); err != nil {
			return report, err
		}
	}
	return report, nil
}

func runBaselineCase(options CorpusOptions, item CorpusCase, sourcePath, caseDir string, result CorpusCaseResult) CorpusCaseResult {
	baselineDir := filepath.Join(caseDir, "baseline")
	candidateDir := filepath.Join(caseDir, "candidate")
	if err := os.MkdirAll(baselineDir, 0o755); err != nil {
		result.TerminalResult, result.Error = "FAIL", err.Error()
		return result
	}
	started := time.Now()
	first := runProcess(options.Compiler, []string{"--phase", options.Phase, "--input", sourcePath, "--input-kind", "source", "--source", sourcePath, "--output-dir", baselineDir, "--run-id", item.ID + "-baseline", "--role", "baseline"}, filepath.Join(caseDir, "baseline.stdout"), filepath.Join(caseDir, "baseline.stderr"))
	result.WallMS += first.WallMS
	result.PeakRSSKiB = first.PeakRSSKiB
	if first.Status != 0 {
		result.TerminalResult, result.Error = "FAIL", first.Error
		return result
	}
	if err := os.MkdirAll(candidateDir, 0o755); err != nil {
		result.TerminalResult, result.Error = "FAIL", err.Error()
		return result
	}
	second := runProcess(options.Compiler, []string{"--phase", options.Phase, "--input", filepath.Join(baselineDir, "semantic-ir.json"), "--input-kind", "semantic-ir", "--source", sourcePath, "--output-dir", candidateDir, "--run-id", item.ID + "-candidate", "--role", "candidate"}, filepath.Join(caseDir, "candidate.stdout"), filepath.Join(caseDir, "candidate.stderr"))
	result.WallMS += second.WallMS
	if second.PeakRSSKiB > result.PeakRSSKiB {
		result.PeakRSSKiB = second.PeakRSSKiB
	}
	if second.Status != 0 {
		result.TerminalResult, result.Error = "FAIL", second.Error
		return result
	}
	verificationPath := filepath.Join(caseDir, "independent-verification.json")
	third := runProcess(options.Verifier, []string{"--phase", options.Phase, "--source", sourcePath, "--baseline-dir", baselineDir, "--candidate-dir", candidateDir, "--expected", item.Expected, "--output", verificationPath}, filepath.Join(caseDir, "verify.stdout"), filepath.Join(caseDir, "verify.stderr"))
	result.WallMS += third.WallMS
	if third.PeakRSSKiB > result.PeakRSSKiB {
		result.PeakRSSKiB = third.PeakRSSKiB
	}
	if third.Status != 0 {
		result.TerminalResult, result.Error = "FAIL", third.Error
		return result
	}
	var verification struct {
		Decision             string `json:"decision"`
		FirstExecutionDigest string `json:"first_execution_digest"`
		RerunExecutionDigest string `json:"rerun_execution_digest"`
	}
	if err := ReadJSON(verificationPath, &verification); err != nil {
		result.TerminalResult, result.Error = "FAIL", err.Error()
		return result
	}
	result.Observed = verification.Decision
	result.TerminalResult = "PASS"
	result.ReplayDigestMatch = verification.FirstExecutionDigest != "" && verification.FirstExecutionDigest == verification.RerunExecutionDigest
	result.GeneratedArtifacts = true
	if !result.ReplayDigestMatch {
		result.Error = "released independent verifier did not observe equal replay digests"
	}
	_ = started
	return result
}

func runCandidateCase(options CorpusOptions, item CorpusCase, sourcePath, caseDir string, result CorpusCaseResult) CorpusCaseResult {
	started := time.Now()
	observation := runProcess(options.Compiler, []string{"--phase", options.Phase, "--input", sourcePath, "--input-kind", "source", "--source", sourcePath, "--output-dir", filepath.Join(caseDir, "candidate"), "--run-id", item.ID + "-candidate", "--role", "candidate"}, filepath.Join(caseDir, "candidate.stdout"), filepath.Join(caseDir, "candidate.stderr"))
	result.WallMS = observation.WallMS
	result.PeakRSSKiB = observation.PeakRSSKiB
	if observation.Status == 0 {
		var receipt struct {
			Decision string `json:"decision"`
		}
		if err := ReadJSON(filepath.Join(caseDir, "candidate", "receipt.json"), &receipt); err == nil {
			result.Observed = receipt.Decision
			result.GeneratedArtifacts = true
			result.TerminalResult = "PASS"
			result.ReplayDigestMatch = false
			result.Error = "candidate unexpectedly accepted by released compiler interface"
			return result
		}
	}
	result.Observed = "REFUTED"
	result.TerminalResult = "FAIL"
	result.Error = observation.Error
	result.ReplayDigestMatch = false
	_ = started
	return result
}

func runProcess(program string, args []string, stdoutPath, stderrPath string) processObservation {
	started := time.Now()
	stdout, err := os.Create(stdoutPath)
	if err != nil {
		return processObservation{Status: 1, Error: err.Error()}
	}
	defer stdout.Close()
	stderr, err := os.Create(stderrPath)
	if err != nil {
		return processObservation{Status: 1, Error: err.Error()}
	}
	defer stderr.Close()
	cmd := exec.Command(program, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	result := processObservation{WallMS: time.Since(started).Milliseconds()}
	if result.WallMS < 1 {
		result.WallMS = 1
	}
	if state := cmd.ProcessState; state != nil {
		if usage, ok := state.SysUsage().(*syscall.Rusage); ok {
			result.PeakRSSKiB = usage.Maxrss
		}
	}
	if err != nil {
		result.Status = 1
		result.Error = err.Error()
		if data, readErr := os.ReadFile(stderrPath); readErr == nil && strings.TrimSpace(string(data)) != "" {
			result.Error = strings.TrimSpace(string(data))
		}
	}
	return result
}
