package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kimjooyoon/gooo-evolution-trial/internal/trial"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "authority-verify":
		err = authorityVerify(os.Args[2:])
	case "prepare-input":
		err = prepareInput(os.Args[2:])
	case "generate-candidate-phase":
		err = generateCandidatePhase(os.Args[2:])
	case "run-corpus":
		err = runCorpus(os.Args[2:])
	case "prepare-causal-case":
		err = prepareCausalCase(os.Args[2:])
	case "dossier":
		err = dossier(os.Args[2:])
	case "inventory":
		err = inventory(os.Args[2:])
	case "help", "--help", "-h":
		usage()
		return
	default:
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func authorityVerify(args []string) error {
	flags := flag.NewFlagSet("authority-verify", flag.ContinueOnError)
	meta := flags.String("meta", "meta/evolution-trial.gooo", "authority .gooo path")
	contract := flags.String("contract", "contracts/evolution-trial-denominator-v1.json", "denominator path")
	output := flags.String("output", "", "optional output path")
	if err := flags.Parse(args); err != nil {
		return err
	}
	report, err := trial.VerifyAuthority(*meta, *contract)
	if err != nil {
		return err
	}
	if *output != "" {
		if err := trial.WriteJSON(*output, report); err != nil {
			return err
		}
	}
	if err := printJSON(report); err != nil {
		return err
	}
	if report.Decision != "CLOSED" {
		return fmt.Errorf("authority decision is %s", report.Decision)
	}
	return nil
}

func prepareInput(args []string) error {
	flags := flag.NewFlagSet("prepare-input", flag.ContinueOnError)
	phase := flags.String("phase", "", "released compiler phase path")
	commit := flags.String("compiler-commit", "", "immutable compiler release target commit")
	baselineReceipt := flags.String("baseline-receipt", "", "optional real baseline receipt path")
	output := flags.String("output-dir", "", "empty caller-owned output directory")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *phase == "" || *commit == "" || *output == "" {
		return fmt.Errorf("phase, compiler-commit, and output-dir are required")
	}
	return trial.PrepareDeltaInput(*phase, *commit, *baselineReceipt, *output)
}

func generateCandidatePhase(args []string) error {
	flags := flag.NewFlagSet("generate-candidate-phase", flag.ContinueOnError)
	candidate := flags.String("candidate", "", "released delta-forge candidate bundle")
	output := flags.String("output", "", "caller-owned candidate .gooo path")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *candidate == "" || *output == "" {
		return fmt.Errorf("candidate and output are required")
	}
	return trial.GenerateCandidatePhase(*candidate, *output)
}

func runCorpus(args []string) error {
	flags := flag.NewFlagSet("run-corpus", flag.ContinueOnError)
	compiler := flags.String("compiler", "", "released compiler executable")
	verifier := flags.String("verifier", "", "released independent verifier executable")
	phase := flags.String("phase", "", "phase .gooo path")
	root := flags.String("root", "", "released compiler source root")
	role := flags.String("role", "", "baseline or candidate")
	output := flags.String("output-dir", "", "empty caller-owned output directory")
	report := flags.String("report", "", "corpus report path")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *compiler == "" || *phase == "" || *root == "" || *role == "" || *output == "" || *report == "" {
		return fmt.Errorf("compiler, phase, root, role, output-dir, and report are required")
	}
	result, err := trial.RunCorpus(trial.CorpusOptions{Compiler: *compiler, Verifier: *verifier, Phase: *phase, Root: *root, Role: *role, OutputDir: *output, ReportPath: *report})
	if err != nil {
		return err
	}
	return printJSON(result)
}

func prepareCausalCase(args []string) error {
	flags := flag.NewFlagSet("prepare-causal-case", flag.ContinueOnError)
	baseline := flags.String("baseline-report", "", "baseline corpus report")
	candidate := flags.String("candidate-report", "", "candidate corpus report")
	releaseID := flags.Int64("causal-release-id", 0, "immutable causal runner release ID")
	assetName := flags.String("causal-asset-name", "", "immutable causal runner release asset name")
	assetDigest := flags.String("causal-asset-digest", "", "immutable causal runner release asset digest")
	output := flags.String("output", "", "caller-owned causal case path")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *baseline == "" || *candidate == "" || *releaseID == 0 || *assetName == "" || *assetDigest == "" || *output == "" {
		return fmt.Errorf("baseline-report, candidate-report, causal release coordinates, and output are required")
	}
	return trial.PrepareCausalCase(*baseline, *candidate, *output, *releaseID, *assetName, *assetDigest)
}

func dossier(args []string) error {
	flags := flag.NewFlagSet("dossier", flag.ContinueOnError)
	baseline := flags.String("baseline-report", "", "baseline corpus report")
	candidate := flags.String("candidate-report", "", "candidate corpus report")
	bundle := flags.String("candidate-bundle", "", "delta-forge candidate bundle")
	evidence := flags.String("normalization-evidence", "", "normalization evidence")
	causal := flags.String("causal-receipt", "", "causal runner verification receipt")
	process := flags.String("process-evidence", "", "GitHub REST process evidence")
	metrics := flags.String("metrics", "", "exact CI metrics receipt")
	report := flags.String("report", "", "final machine report")
	dossierPath := flags.String("dossier", "", "human dossier")
	if err := flags.Parse(args); err != nil {
		return err
	}
	for _, value := range []*string{baseline, candidate, bundle, evidence, causal, process, metrics, report, dossierPath} {
		if *value == "" {
			return fmt.Errorf("all dossier paths are required")
		}
	}
	return trial.WriteDossier(*baseline, *candidate, *bundle, *evidence, *causal, *process, *metrics, *report, *dossierPath)
}

func inventory(args []string) error {
	flags := flag.NewFlagSet("inventory", flag.ContinueOnError)
	root := flags.String("root", ".", "repository root")
	if err := flags.Parse(args); err != nil {
		return err
	}
	type inventoryReport struct {
		Schema             string `json:"schema"`
		RegularFiles       int    `json:"regular_files"`
		Subdirectories     int    `json:"subdirectories"`
		GoFiles            int    `json:"go_files"`
		GoooFiles          int    `json:"gooo_files"`
		PhysicalLines      int    `json:"physical_lines"`
		GoLines            int    `json:"go_lines"`
		GoooLines          int    `json:"gooo_lines"`
		RootReadmeExcluded bool   `json:"root_readme_excluded"`
	}
	result := inventoryReport{Schema: "gooo/evolution-trial/inventory/v1", RootReadmeExcluded: true}
	err := filepath.WalkDir(*root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			if path != *root {
				result.Subdirectories++
			}
			return nil
		}
		if !entry.Type().IsRegular() || path == filepath.Join(*root, "README.md") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		result.RegularFiles++
		result.PhysicalLines += lineCount(data)
		switch filepath.Ext(path) {
		case ".go":
			result.GoFiles++
			result.GoLines += lineCount(data)
		case ".gooo":
			result.GoooFiles++
			result.GoooLines += lineCount(data)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return printJSON(result)
}

func lineCount(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	count := 1
	for _, value := range data {
		if value == '\n' {
			count++
		}
	}
	if data[len(data)-1] == '\n' {
		count--
	}
	return count
}

func printJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func usage() {
	fmt.Fprintln(os.Stderr, "gooo-evolution-trial authority-verify|prepare-input|generate-candidate-phase|run-corpus|prepare-causal-case|dossier|inventory")
}
