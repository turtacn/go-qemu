package harness

import (
	"context"
	"fmt"

	"github.com/turtacn/go-qemu/harness/scenarios"
)

type SuiteResult struct {
	Total  int
	Passed int
	Failed int
	Errors []string
}

type Suite struct {
	scenarios []scenarios.Scenario
}

func NewSuite() *Suite {
	return &Suite{
		scenarios: []scenarios.Scenario{
			scenarios.NewLocalUpgradeScenario(100),
			&scenarios.CrossVersionMigrateScenario{},
			&scenarios.CmdlineDiffScenario{},
			&scenarios.Cfg2CmdCompatScenario{},
			&scenarios.FeatureMatrixScenario{},
		},
	}
}

func (s *Suite) Run(ctx context.Context) (*SuiteResult, error) {
	res := &SuiteResult{
		Total: len(s.scenarios),
	}

	for _, sc := range s.scenarios {
		fmt.Printf("Running Scenario: %s (%s)...\n", sc.Name(), sc.Description())
		if err := sc.Prerequisites(ctx); err != nil {
			errStr := fmt.Sprintf("Prerequisites failed for %s: %v", sc.Name(), err)
			res.Failed++
			res.Errors = append(res.Errors, errStr)
			continue
		}
		if err := sc.Setup(ctx); err != nil {
			errStr := fmt.Sprintf("Setup failed for %s: %v", sc.Name(), err)
			res.Failed++
			res.Errors = append(res.Errors, errStr)
			continue
		}
		execErr := sc.Execute(ctx)
		verifyErr := sc.Verify(ctx)
		_ = sc.Teardown(ctx)

		if execErr != nil {
			errStr := fmt.Sprintf("Execute failed for %s: %v", sc.Name(), execErr)
			res.Failed++
			res.Errors = append(res.Errors, errStr)
			continue
		}
		if verifyErr != nil {
			errStr := fmt.Sprintf("Verify failed for %s: %v", sc.Name(), verifyErr)
			res.Failed++
			res.Errors = append(res.Errors, errStr)
			continue
		}

		res.Passed++
		fmt.Printf("Scenario %s PASSED.\n", sc.Name())
	}

	if res.Failed > 0 {
		return res, fmt.Errorf("suite completed with %d failed scenario(s)", res.Failed)
	}
	return res, nil
}
