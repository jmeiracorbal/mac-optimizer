package modules

func FreeMemory(ctx *Ctx) error {
	ctx.Log.Progress("Freeing system memory")

	if err := ctx.Runner.Sudo("purge"); err != nil {
		return err
	}

	ctx.Log.Success("Memory purge complete")
	return nil
}
