package matrix

import (
	"context"
	"os"
	"testing"

	"github.com/turtacn/go-qemu/harness"
	"gopkg.in/yaml.v3"
)

type MatrixConfig struct {
	Matrix struct {
		QEMUVersions  []string `yaml:"qemu_versions"`
		Architectures []string `yaml:"architectures"`
	} `yaml:"matrix"`
}

func TestMatrixConfig(t *testing.T) {
	data, err := os.ReadFile("matrix.yaml")
	if err != nil {
		t.Fatalf("failed to read matrix.yaml: %v", err)
	}

	var cfg MatrixConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to unmarshal matrix.yaml: %v", err)
	}

	if len(cfg.Matrix.QEMUVersions) < 2 {
		t.Errorf("expected at least 2 QEMU versions in matrix.yaml")
	}
}

func TestRunSuite(t *testing.T) {
	suite := harness.NewSuite()
	res, err := suite.Run(context.Background())
	if err != nil {
		t.Fatalf("Suite run failed: %v", err)
	}
	if res.Passed != res.Total {
		t.Errorf("expected all scenarios to pass, passed %d / %d", res.Passed, res.Total)
	}
}
