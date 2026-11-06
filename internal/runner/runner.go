package runner

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/jmeiracorbal/mac-optimizer/internal/logger"
)

type Runner struct {
	dryRun bool
	log    *logger.Logger
}

func New(dryRun bool, log *logger.Logger) *Runner {
	return &Runner{dryRun: dryRun, log: log}
}

func (r *Runner) IsDryRun() bool { return r.dryRun }

// Exec runs a command. No-ops in dry-run mode.
func (r *Runner) Exec(name string, args ...string) error {
	if r.dryRun {
		r.log.Info("[dry-run] %s", cmdLine(name, args))
		return nil
	}
	if err := exec.Command(name, args...).Run(); err != nil {
		return fmt.Errorf("%s: %w", cmdLine(name, args), err)
	}
	return nil
}

// Sudo runs args under sudo. No-ops in dry-run mode.
func (r *Runner) Sudo(args ...string) error {
	return r.Exec("sudo", args...)
}

// Query runs a read-only command and captures its stdout.
// Always executes, even in dry-run mode.
func (r *Runner) Query(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

func cmdLine(name string, args []string) string {
	return strings.Join(append([]string{name}, args...), " ")
}
