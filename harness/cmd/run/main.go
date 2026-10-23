package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/turtacn/go-qemu/harness"
)

func main() {
	matrixFile := flag.String("matrix", "harness/matrix/matrix.yaml", "Path to matrix.yaml")
	reportFile := flag.String("report", "/tmp/harness-report.json", "Path to output report JSON")
	flag.Parse()

	fmt.Printf("Running Harness SOP Matrix using %s...\n", *matrixFile)

	suite := harness.NewSuite()
	ctx := context.Background()
	res, err := suite.Run(ctx)

	status := "PASSED"
	if err != nil {
		status = "FAILED"
		log.Printf("Harness SOP execution failed: %v", err)
	}

	report := &harness.TestReport{
		TotalScenarios:  res.Total,
		PassedScenarios: res.Passed,
		FailedScenarios: res.Failed,
		Status:          status,
		Details:         res.Errors,
	}

	if err := harness.SaveReport(report, *reportFile); err != nil {
		log.Printf("Failed to save report to %s: %v", *reportFile, err)
	} else {
		fmt.Printf("Harness report saved to %s (Passed: %d/%d)\n", *reportFile, res.Passed, res.Total)
	}
}
