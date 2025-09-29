package main

import (
	"log"
	fxmodules "pets_rest/internal/fx"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	fx.New(
		// core modules
		fxmodules.ConfigModule,
		fxmodules.DatabaseModule,
		// storage module
		fxmodules.StorageModule,
		//repositories module
		fxmodules.RepositoryModule,
		// business logic module
		fxmodules.ServiceModule,
		// http layer module
		fxmodules.FiberModule,
		fxmodules.HandlerModule,
		fxmodules.RouteModule,

		// server lifecycle module
		fxmodules.ServerModule,

		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ConsoleLogger{W: log.Writer()}
		}),
	).Run()

}
