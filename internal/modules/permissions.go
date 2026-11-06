package modules

import (
	"os"
)

func RepairPermissions(ctx *Ctx) error {
	ctx.Log.Progress("Repairing system permissions")

	// diskutil repairPermissions was removed in macOS Sierra; the call is a no-op
	// on Monterey+ but we keep it for older versions that may still respond.
	if err := ctx.Runner.Sudo("diskutil", "repairPermissions", "/"); err != nil {
		ctx.Log.Info("diskutil repairPermissions not available on this macOS version")
	}

	if err := ctx.Runner.Sudo("chmod", "755", "/usr/local/bin"); err != nil {
		ctx.Log.Warn("Could not set permissions on /usr/local/bin: %v", err)
	}

	home, _ := os.UserHomeDir()
	if !ctx.DryRun {
		if err := os.Chmod(home, 0755); err != nil {
			ctx.Log.Warn("Could not set permissions on %s: %v", home, err)
		}
	} else {
		ctx.Log.Info("[dry-run] would chmod 755 %s", home)
	}

	ctx.Log.Success("Permissions repaired")
	return nil
}
