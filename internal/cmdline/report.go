package cmdline

import (
	"fmt"
	"io"

	"github.com/turtacn/go-qemu/internal/feature"
	"github.com/turtacn/go-qemu/internal/version"
)

type RiskLevel int

const (
	RiskNone RiskLevel = iota
	RiskDeprecated
	RiskBehaviorChange
	RiskBreaking
)

type ArgEntry struct {
	Key         string
	Value       string
	RuleName    string
	FeatureGate string
	Reason      string
	Deprecated  bool
	Risk        RiskLevel
}

func (a ArgEntry) String() string {
	if a.Value == "" {
		return a.Key
	}
	return fmt.Sprintf("%s %s", a.Key, a.Value)
}

type BuildReport struct {
	QEMUVersion version.QEMUVersion
	Arch        feature.Arch
	Args        []ArgEntry
	Warnings    []string
	Errors      []string
}

func (r *BuildReport) ToArgSlice() []string {
	result := make([]string, 0, len(r.Args)*2)
	for _, a := range r.Args {
		result = append(result, a.Key)
		if a.Value != "" {
			result = append(result, a.Value)
		}
	}
	return result
}

func (r *BuildReport) HasRisk(level RiskLevel) bool {
	for _, a := range r.Args {
		if a.Risk >= level {
			return true
		}
	}
	return false
}

type CmdlineDiff struct {
	FromVersion version.QEMUVersion
	ToVersion   version.QEMUVersion
	Added       []ArgEntry
	Removed     []ArgEntry
	Changed     []ArgChange
	Unchanged   []ArgEntry
}

type ArgChange struct {
	From ArgEntry
	To   ArgEntry
}

func (d *CmdlineDiff) HasRisk(level RiskLevel) bool {
	for _, a := range d.Added {
		if a.Risk >= level {
			return true
		}
	}
	for _, c := range d.Changed {
		if c.To.Risk >= level {
			return true
		}
	}
	return false
}

func (d *CmdlineDiff) Print(w io.Writer) {
	fmt.Fprintf(w, "=== 命令行版本差异报告 ===\n")
	fmt.Fprintf(w, "源版本: %s → 目标版本: %s\n\n", d.FromVersion, d.ToVersion)

	if len(d.Added) > 0 {
		fmt.Fprintf(w, "-- 新增参数 --\n")
		for _, a := range d.Added {
			fmt.Fprintf(w, "+ %s %s\n  规则: %s | 理由: %s\n", a.Key, a.Value, a.RuleName, a.Reason)
		}
	}
	if len(d.Changed) > 0 {
		fmt.Fprintf(w, "\n-- 参数变更 --\n")
		for _, c := range d.Changed {
			risk := ""
			if c.To.Risk >= RiskBehaviorChange {
				risk = " [行为变更警告]"
			}
			fmt.Fprintf(w, "~ %s %s  ->  %s%s\n  理由: %s\n",
				c.From.Key, c.From.Value, c.To.Value, risk, c.To.Reason)
		}
	}
	if len(d.Removed) > 0 {
		fmt.Fprintf(w, "\n-- 删除参数 --\n")
		for _, a := range d.Removed {
			fmt.Fprintf(w, "- %s %s\n", a.Key, a.Value)
		}
	}
	fmt.Fprintf(w, "\n无变化参数: %d 个\n", len(d.Unchanged))
}
