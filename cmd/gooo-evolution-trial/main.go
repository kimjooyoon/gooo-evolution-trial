package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 || os.Args[1] != "authority-verify" {
		fmt.Fprintln(os.Stderr, "usage: gooo-evolution-trial authority-verify")
		os.Exit(2)
	}
	data, err := os.ReadFile("meta/evolution-trial.gooo")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	text := string(data)
	for _, marker := range []string{
		"gooo evolution_trial v1",
		"precedence REFUTED>UNKNOWN>CLOSED",
		"unknown_fields stage,step,reason,unknown_class,next_operation,blocked_by",
		"authority repository_writes=0 upstream_writes=0 local_test_executions=0 merge_authority=0",
		"activity AdoptOrRollback",
		"process_field transition \"PR_MERGE_ONLY\"",
	} {
		if !strings.Contains(text, marker) {
			fmt.Fprintf(os.Stderr, "authority marker missing: %s\n", marker)
			os.Exit(1)
		}
	}
	fmt.Println(`{"decision":"CLOSED","repository_writes":0,"upstream_writes":0,"local_test_executions":0}`)
}
