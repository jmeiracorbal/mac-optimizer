package modules

import (
	"os"
	"path/filepath"
)

func CleanBrowsers(ctx *Ctx) error {
	ctx.Log.Progress("Cleaning browser caches")

	if !ctx.Confirm("Browser cleanup clears caches and local storage — you may need to log in to some sites again") {
		ctx.Log.Info("Browser cleanup skipped")
		return nil
	}

	home, _ := os.UserHomeDir()

	targets := []string{
		// Chrome
		filepath.Join(home, "Library", "Caches", "Google", "Chrome"),
		filepath.Join(home, "Library", "Application Support", "Google", "Chrome", "Default", "Application Cache"),
		filepath.Join(home, "Library", "Application Support", "Google", "Chrome", "Default", "GPUCache"),
		// Firefox
		filepath.Join(home, "Library", "Application Support", "Firefox", "Profiles"),
		// Safari
		filepath.Join(home, "Library", "Caches", "com.apple.Safari"),
		filepath.Join(home, "Library", "Safari", "LocalStorage"),
		filepath.Join(home, "Library", "Safari", "Databases"),
		filepath.Join(home, "Library", "WebKit"),
	}

	for _, dir := range targets {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		before := dirSize(dir)
		ctx.Log.Info("Cleaning %s (%s)", dir, formatSize(before))
		freed := emptyDir(ctx, dir)
		if !ctx.DryRun && freed > 0 {
			ctx.Log.Success("Freed %s from %s", formatSize(freed), filepath.Base(dir))
		}
	}

	ctx.Log.Success("Browser cache cleanup complete")
	return nil
}
