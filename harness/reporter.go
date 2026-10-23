package harness

import (
	"encoding/json"
	"os"
)

type TestReport struct {
	TotalScenarios  int      `json:"total_scenarios"`
	PassedScenarios int      `json:"passed_scenarios"`
	FailedScenarios int      `json:"failed_scenarios"`
	Status          string   `json:"status"`
	Details         []string `json:"details"`
}

func SaveReport(report *TestReport, path string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
