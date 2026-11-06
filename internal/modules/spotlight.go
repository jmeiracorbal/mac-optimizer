package modules

func OptimizeSpotlight(ctx *Ctx) error {
	ctx.Log.Progress("Spotlight optimization requested")

	if !ctx.Confirm("Spotlight optimization requires admin permissions and rebuilds the full index (may take time)") {
		ctx.Log.Info("Spotlight optimization skipped")
		return nil
	}

	ctx.Log.Info("Disabling Spotlight indexing")
	if err := ctx.Runner.Sudo("mdutil", "-a", "-i", "off"); err != nil {
		return err
	}

	ctx.Log.Info("Removing existing Spotlight index")
	ctx.Runner.Sudo("rm", "-rf", "/.Spotlight-V100", "/System/Volumes/Data/.Spotlight-V100")

	ctx.Log.Info("Re-enabling Spotlight indexing")
	if err := ctx.Runner.Sudo("mdutil", "-a", "-i", "on"); err != nil {
		return err
	}

	ctx.Log.Success("Spotlight reindex started (continues in background)")
	return nil
}
