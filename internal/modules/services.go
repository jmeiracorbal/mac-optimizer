package modules

func RestartServices(ctx *Ctx) error {
	ctx.Log.Progress("Restarting system services")

	services := []string{"Finder", "Dock", "SystemUIServer", "ControlStrip"}

	for _, svc := range services {
		ctx.Log.Info("Restarting %s", svc)
		if err := ctx.Runner.Exec("killall", svc); err != nil {
			ctx.Log.Warn("%s not running or could not be restarted", svc)
		}
	}

	ctx.Log.Success("Services restarted")
	return nil
}
