package modules

import (
	"os"
	"path/filepath"
)

func CleanLogs(ctx *Ctx) error {
	ctx.Log.Progress("Cleaning system logs")

	home, _ := os.UserHomeDir()

	dirs := []string{
		filepath.Join(home, "Library", "Logs"),
		"/private/var/log",
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		ctx.Log.Info("Cleaning logs in %s", dir)
		freed := deleteOldFiles(ctx, dir, 30, ".log", ".out")
		if !ctx.DryRun && freed > 0 {
			ctx.Log.Info("Freed %s in %s", formatSize(freed), dir)
		}
	}

	ctx.Log.Success("Log cleanup complete")
	return nil
}
