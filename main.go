package main

import (
	"DistanceTrackerServer/constants"
	"DistanceTrackerServer/router"
	"DistanceTrackerServer/utils"
	"fmt"
	"go.uber.org/zap"
	"time"
)

var (
	initRouter = router.Init
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

	sugar.Infof("Starting server on port %s", constants.ServerPort)
	if err := app.Run(fmt.Sprintf(":%s", constants.ServerPort)); err != nil {
		sugar.Fatal("Failed to start server: ", err)
	}
}

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	run()
}
