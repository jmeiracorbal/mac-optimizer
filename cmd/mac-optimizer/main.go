package main

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/jmeiracorbal/mac-optimizer/internal/config"
	"github.com/jmeiracorbal/mac-optimizer/internal/logger"
	"github.com/jmeiracorbal/mac-optimizer/internal/modules"
	"github.com/jmeiracorbal/mac-optimizer/internal/runner"
	"github.com/jmeiracorbal/mac-optimizer/internal/stats"
	"github.com/spf13/cobra"
)

// version is overridden at link time by goreleaser.
var version = "dev"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cfg := config.Default()

	cmd := &cobra.Command{
		Use:          "mac-optimizer [cache|spotlight|permissions|services|memory|logs|browsers|downloads|all...]",
		Short:        "macOS optimization and cleanup tool",
		Version:      version,
		Args:         cobra.ArbitraryArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cfg, args)
		},
	}

	cmd.SetVersionTemplate("mac-optimizer {{.Version}}\nCompatible with macOS Monterey+\n")

	cmd.Flags().BoolVarP(&cfg.DryRun, "dry-run", "d", false, "Show what would be done without making changes")
	cmd.Flags().BoolVarP(&cfg.Quiet, "quiet", "q", false, "Run in quiet mode")
	cmd.Flags().BoolVarP(&cfg.Force, "force", "f", false, "Force operations without confirmation")
	cmd.Flags().StringVar(&cfg.LogFile, "log-file", cfg.LogFile, "Log file path")

	return cmd
}

func run(cfg *config.Config, args []string) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("mac-optimizer only supports macOS")
	}

	log, err := logger.New(cfg.LogFile, cfg.Quiet)
	if err != nil {
		return fmt.Errorf("open log: %w", err)
	}
	defer log.Close()

	r := runner.New(cfg.DryRun, log)
	ctx := modules.NewCtx(r, log, cfg.Force, cfg.DryRun)

	names, err := resolveModules(args)
	if err != nil {
		return err
	}

	printHeader(cfg, version, r)
	log.Info("=== mac-optimizer %s ===", version)

	for _, name := range names {
		if err := modules.Registry[name](ctx); err != nil {
			log.Error("module %s: %v", name, err)
		}
		fmt.Println()
	}

	fmt.Printf("\033[0;32m✓  Optimization complete. Log: %s\033[0m\n\n", cfg.LogFile)
	stats.Print(r)

	return nil
}

func resolveModules(args []string) ([]string, error) {
	if len(args) == 0 {
		return modules.DefaultOrder, nil
	}

	var names []string
	for _, a := range args {
		if a == "all" {
			names = append(names, modules.DefaultOrder...)
			continue
		}
		if _, ok := modules.Registry[a]; !ok {
			available := modules.AvailableNames()
			sort.Strings(available)
			return nil, fmt.Errorf("unknown module %q — available: %s", a, strings.Join(available, ", "))
		}
		names = append(names, a)
	}
	return names, nil
}

func printHeader(cfg *config.Config, ver string, r *runner.Runner) {
	fmt.Printf("\033[0;36m=== mac-optimizer %s ===\033[0m\n", ver)
	if cfg.DryRun {
		fmt.Println("\033[1;33m[dry-run] showing what would be done\033[0m")
	}
	fmt.Println()
	stats.Print(r)
}
