package main

import (
	"DistanceTrackerServer/router"
	"DistanceTrackerServer/utils"
	"go.uber.org/zap"
	"time"
)

var (
	sugarFromContext = utils.SugarFromContext
	initRouter       = router.Init
)

func heartbeat(sugar *zap.SugaredLogger) {
	for {
		time.Sleep(60 * time.Second)
		sugar.Info("HEARTBEAT")
	}
}

func run() {
	log := utils.Logger
	sugar := utils.Sugar
	sugar.Info("Starting application!!")

	go heartbeat(sugar)

	app := initRouter(log)
	defer func() {
		err := sugar.Sync()
		if err != nil {
			sugar.Error("Failed to sync logger: ", err)
		}
	}()

	sugar.Info("Starting server on port 8080")
	if err := app.Run(":8080"); err != nil {
		sugar.Fatal("Failed to start server: ", err)
	}
}

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	run()
}
