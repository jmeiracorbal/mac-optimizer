package modules

import (
	"os"
	"path/filepath"
)

func CleanCaches(ctx *Ctx) error {
	ctx.Log.Progress("Starting cache cleanup")

	home, _ := os.UserHomeDir()

	userDirs := []string{
		filepath.Join(home, "Library", "Caches"),
		filepath.Join(home, "Library", "Application Support", "CrashReporter"),
		filepath.Join(home, "Library", "Logs"),
		filepath.Join(home, "Library", "WebKit"),
		filepath.Join(home, "Library", "Saved Application State"),
	}

	for _, dir := range userDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		before := dirSize(dir)
		ctx.Log.Info("Cleaning %s (%s)", dir, formatSize(before))
		freed := emptyDir(ctx, dir)
		if !ctx.DryRun {
			ctx.Log.Success("Cleaned %s (freed %s)", dir, formatSize(freed))
		}
	}

	systemDirs := []string{
		"/Library/Caches",
		"/private/var/log/asl",
	}

	for _, dir := range systemDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		ctx.Log.Info("Cleaning system cache: %s", dir)
		if err := ctx.Runner.Sudo("find", dir, "-type", "f", "-mtime", "+1", "-delete"); err != nil {
			ctx.Log.Warn("Could not clean %s: %v", dir, err)
		}
	}

	ctx.Log.Info("Flushing DNS cache")
	ctx.Runner.Sudo("dscacheutil", "-flushcache")
	ctx.Runner.Sudo("killall", "-HUP", "mDNSResponder")

	ctx.Log.Success("Cache cleanup complete")
	return nil
}
