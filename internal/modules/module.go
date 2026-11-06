package modules

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmeiracorbal/mac-optimizer/internal/logger"
	"github.com/jmeiracorbal/mac-optimizer/internal/runner"
)

// ModFunc is the signature every optimization module implements.
type ModFunc func(*Ctx) error

// Ctx carries shared state into every module.
type Ctx struct {
	Runner *runner.Runner
	Log    *logger.Logger
	Force  bool
	DryRun bool
}

func NewCtx(r *runner.Runner, log *logger.Logger, force, dryRun bool) *Ctx {
	return &Ctx{Runner: r, Log: log, Force: force, DryRun: dryRun}
}

// Confirm prompts the user unless Force or DryRun is set.
func (ctx *Ctx) Confirm(msg string) bool {
	if ctx.Force || ctx.DryRun {
		return true
	}
	fmt.Printf("\033[1;33m⚠  %s (y/N): \033[0m", msg)
	var input string
	fmt.Scanln(&input)
	return strings.EqualFold(strings.TrimSpace(input), "y")
}

// emptyDir deletes all files inside path. Returns bytes freed.
func emptyDir(ctx *Ctx, path string) int64 {
	if ctx.DryRun {
		size := dirSize(path)
		ctx.Log.Info("[dry-run] would empty %s (%s)", path, formatSize(size))
		return size
	}
	var freed int64
	filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err == nil {
			freed += info.Size()
		}
		os.Remove(p)
		return nil
	})
	return freed
}

// deleteOldFiles removes files older than maxAge days under dir, filtered by extensions
// (empty extensions list matches all files). Returns bytes freed.
func deleteOldFiles(ctx *Ctx, dir string, maxAge int, extensions ...string) int64 {
	cutoff := time.Now().AddDate(0, 0, -maxAge)
	var freed int64
	filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if len(extensions) > 0 && !matchesExt(p, extensions) {
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.ModTime().Before(cutoff) {
			return nil
		}
		freed += info.Size()
		if !ctx.DryRun {
			os.Remove(p)
		}
		return nil
	})
	return freed
}

func dirSize(path string) int64 {
	var size int64
	filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err == nil {
			size += info.Size()
		}
		return nil
	})
	return size
}

func countOldFiles(dir string, maxAge int) (int, int64) {
	cutoff := time.Now().AddDate(0, 0, -maxAge)
	var count int
	var size int64
	filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.ModTime().Before(cutoff) {
			return nil
		}
		count++
		size += info.Size()
		return nil
	})
	return count, size
}

func matchesExt(path string, exts []string) bool {
	ext := filepath.Ext(path)
	for _, e := range exts {
		if ext == e {
			return true
		}
	}
	return false
}

func formatSize(b int64) string {
	if b == 0 {
		return "0 B"
	}
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// Registry maps module names to their implementation.
var Registry = map[string]ModFunc{
	"cache":       CleanCaches,
	"spotlight":   OptimizeSpotlight,
	"permissions": RepairPermissions,
	"services":    RestartServices,
	"memory":      FreeMemory,
	"logs":        CleanLogs,
	"browsers":    CleanBrowsers,
	"downloads":   CleanDownloads,
}

// DefaultOrder is the module sequence when none are specified.
var DefaultOrder = []string{
	"cache", "spotlight", "permissions", "services", "memory", "logs",
}

// AvailableNames returns all registered module names.
func AvailableNames() []string {
	names := make([]string, 0, len(Registry))
	for k := range Registry {
		names = append(names, k)
	}
	return names
}
