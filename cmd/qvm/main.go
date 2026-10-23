package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/turtacn/go-qemu/harness"
	"github.com/turtacn/go-qemu/internal/cmdline"
	"github.com/turtacn/go-qemu/internal/cmdline/rules"
	"github.com/turtacn/go-qemu/internal/dryrun"
	"github.com/turtacn/go-qemu/internal/upgrade"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcmd := os.Args[1]
	ctx := context.Background()

	switch subcmd {
	case "build-cmdline":
		fs := flag.NewFlagSet("build-cmdline", flag.ExitOnError)
		cfgPath := fs.String("config", "", "Path to VM configuration YAML")
		qemuVer := fs.String("version", "5.0.0", "Target QEMU version")
		fs.Parse(os.Args[2:])

		if *cfgPath == "" {
			fs.Usage()
			os.Exit(1)
		}

		vmCfg, err := vmconfig.LoadFromHCI(*cfgPath)
		if err != nil {
			log.Fatalf("Failed to load VM config: %v", err)
		}

		parsedVer, err := version.Parse(*qemuVer)
		if err != nil {
			log.Fatalf("Invalid QEMU version: %v", err)
		}

		builder, err := cmdline.NewBuilder(cmdline.WithQEMUVersion(parsedVer))
		if err != nil {
			log.Fatalf("Failed to create builder: %v", err)
		}
		reg := rules.DefaultRegistry()
		var ruleSlice []cmdline.Rule
		for _, r := range reg.Rules() {
			ruleSlice = append(ruleSlice, r)
		}
		builder.SetRules(ruleSlice)

		args, report, err := builder.Build(vmCfg)
		if err != nil {
			log.Fatalf("Build failed: %v", err)
		}

		fmt.Printf("Generated command line for QEMU %s:\n", parsedVer)
		for _, a := range args {
			fmt.Printf("  %s\n", a)
		}
		if len(report.Warnings) > 0 {
			fmt.Println("Warnings:")
			for _, w := range report.Warnings {
				fmt.Printf("  - %s\n", w)
			}
		}

	case "diff":
		fs := flag.NewFlagSet("diff", flag.ExitOnError)
		cfgPath := fs.String("config", "", "Path to VM configuration YAML")
		fromVerStr := fs.String("from", "5.0.0", "Source QEMU version")
		toVerStr := fs.String("to", "10.2.2", "Target QEMU version")
		fs.Parse(os.Args[2:])

		if *cfgPath == "" {
			fs.Usage()
			os.Exit(1)
		}

		vmCfg, err := vmconfig.LoadFromHCI(*cfgPath)
		if err != nil {
			log.Fatalf("Failed to load VM config: %v", err)
		}

		fromVer, _ := version.Parse(*fromVerStr)
		toVer, _ := version.Parse(*toVerStr)

		reg := rules.DefaultRegistry()
		var ruleSlice []cmdline.Rule
		for _, r := range reg.Rules() {
			ruleSlice = append(ruleSlice, r)
		}

		bFrom, _ := cmdline.NewBuilder(cmdline.WithQEMUVersion(fromVer))
		bFrom.SetRules(ruleSlice)
		_, reportFrom, _ := bFrom.Build(vmCfg)

		bTo, _ := cmdline.NewBuilder(cmdline.WithQEMUVersion(toVer))
		bTo.SetRules(ruleSlice)
		_, reportTo, _ := bTo.Build(vmCfg)

		diff := reportTo.DiffFrom(fromVer, reportFrom)
		diff.Print(os.Stdout)

	case "cfg2cmd":
		fs := flag.NewFlagSet("cfg2cmd", flag.ExitOnError)
		pvePath := fs.String("pve-config", "", "Path to Proxmox PVE format configuration")
		qemuVer := fs.String("version", "5.0.0", "Target QEMU version")
		fs.Parse(os.Args[2:])

		if *pvePath == "" {
			fs.Usage()
			os.Exit(1)
		}

		vmCfg, err := vmconfig.LoadFromPVE(*pvePath)
		if err != nil {
			log.Fatalf("Failed to load PVE config: %v", err)
		}

		parsedVer, _ := version.Parse(*qemuVer)
		builder, _ := cmdline.NewBuilder(cmdline.WithQEMUVersion(parsedVer))
		reg := rules.DefaultRegistry()
		var ruleSlice []cmdline.Rule
		for _, r := range reg.Rules() {
			ruleSlice = append(ruleSlice, r)
		}
		builder.SetRules(ruleSlice)

		args, _, err := builder.Build(vmCfg)
		if err != nil {
			log.Fatalf("cfg2cmd conversion failed: %v", err)
		}

		fmt.Printf("Converted PVE Config to QEMU %s commandline:\n", parsedVer)
		for _, a := range args {
			fmt.Printf("  %s\n", a)
		}

	case "dry-run":
		fs := flag.NewFlagSet("dry-run", flag.ExitOnError)
		cfgPath := fs.String("config", "", "Path to VM configuration YAML")
		qemuBin := fs.String("qemu-bin", "", "Path to QEMU binary")
		fs.Parse(os.Args[2:])

		if *cfgPath == "" {
			fs.Usage()
			os.Exit(1)
		}

		vmCfg, err := vmconfig.LoadFromHCI(*cfgPath)
		if err != nil {
			log.Fatalf("Failed to load VM config: %v", err)
		}

		builder, _ := cmdline.NewBuilder(cmdline.WithQEMUVersion(version.V(5, 0, 0)))
		reg := rules.DefaultRegistry()
		var ruleSlice []cmdline.Rule
		for _, r := range reg.Rules() {
			ruleSlice = append(ruleSlice, r)
		}
		builder.SetRules(ruleSlice)

		args, _, err := builder.Build(vmCfg)
		if err != nil {
			log.Fatalf("Failed to build args for dry-run: %v", err)
		}

		res, err := dryrun.Run(ctx, dryrun.Config{
			Args:       args,
			QEMUBinary: *qemuBin,
			Layers:     dryrun.Layer1Static,
			VMConfig:   vmCfg,
		})
		if err != nil {
			log.Fatalf("Dry-run failed with error: %v", err)
		}

		if res.OK {
			fmt.Println("Dry-Run Validation PASSED!")
		} else {
			fmt.Println("Dry-Run Validation FAILED!")
			os.Exit(1)
		}

	case "upgrade":
		fs := flag.NewFlagSet("upgrade", flag.ExitOnError)
		cfgPath := fs.String("config", "", "Path to VM configuration YAML")
		vmid := fs.Int("vmid", 100, "VM ID")
		fs.Parse(os.Args[2:])

		if *cfgPath == "" {
			fs.Usage()
			os.Exit(1)
		}

		upgrader := upgrade.NewLocalUpgrader(upgrade.LocalUpgradeConfig{
			VMID:         *vmid,
			VMConfigPath: *cfgPath,
		})

		res, err := upgrader.RunWithResult(ctx)
		if err != nil {
			log.Fatalf("Upgrade failed: %v", err)
		}

		fmt.Printf("VM %d Upgrade Successful: %s -> %s (Downtime: %dms)\n",
			res.VMID, res.SrcVersion, res.DstVersion, res.DowntimeMs)

	case "harness":
		fs := flag.NewFlagSet("harness", flag.ExitOnError)
		reportPath := fs.String("report", "/tmp/harness-report.json", "Output report path")
		fs.Parse(os.Args[2:])

		suite := harness.NewSuite()
		res, err := suite.Run(ctx)

		status := "PASSED"
		if err != nil {
			status = "FAILED"
		}

		report := &harness.TestReport{
			TotalScenarios:  res.Total,
			PassedScenarios: res.Passed,
			FailedScenarios: res.Failed,
			Status:          status,
			Details:         res.Errors,
		}

		if err := harness.SaveReport(report, *reportPath); err != nil {
			log.Fatalf("Failed to save harness report: %v", err)
		}
		fmt.Printf("Harness SOP complete. Report saved to %s (Passed: %d/%d)\n", *reportPath, res.Passed, res.Total)

	default:
		fmt.Printf("Unknown command: %s\n", subcmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: qvm <command> [args]")
	fmt.Println("Commands:")
	fmt.Println("  build-cmdline --config <yaml> [--version <ver>]")
	fmt.Println("  diff          --config <yaml> [--from <ver1>] [--to <ver2>]")
	fmt.Println("  cfg2cmd       --pve-config <path> [--version <ver>]")
	fmt.Println("  dry-run       --config <yaml> [--qemu-bin <path>]")
	fmt.Println("  upgrade       --config <yaml> [--vmid <id>]")
	fmt.Println("  harness       [--report <path>]")
}
