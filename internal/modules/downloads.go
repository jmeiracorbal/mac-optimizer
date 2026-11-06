package modules

import (
	"fmt"
	"os"
	"path/filepath"
)

func CleanDownloads(ctx *Ctx) error {
	ctx.Log.Progress("Cleaning old downloads")

	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, "Downloads")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		ctx.Log.Info("Downloads directory not found")
		return nil
	}

	count, size := countOldFiles(dir, 30)

	if count == 0 {
		ctx.Log.Info("No files older than 30 days in Downloads")
		return nil
	}

	ctx.Log.Info("Found %d files older than 30 days (%s total)", count, formatSize(size))

	if !ctx.Confirm(fmt.Sprintf("Delete %d files older than 30 days from Downloads (%s)?", count, formatSize(size))) {
		ctx.Log.Info("Downloads cleanup skipped")
		return nil
	}

	freed := deleteOldFiles(ctx, dir, 30)
	if !ctx.DryRun {
		ctx.Log.Success("Freed %s from Downloads", formatSize(freed))
	} else {
		ctx.Log.Info("[dry-run] would free %s from Downloads", formatSize(freed))
	}

	return nil
}
