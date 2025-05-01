package main

import (
	"awesomeProject/router"
	"awesomeProject/utils"
	"context"
	"os"
	"time"
)

var (
	createAndSaveLoggers = utils.CreateAndSaveLoggers
	sugarFromContext     = utils.SugarFromContext
	initRouter           = router.Init
)

func heartbeat(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	sugar, _ := sugarFromContext(ctx)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sugar.Info("HEARTBEAT")
		}
	}
}

func run(ctx context.Context, args []string) error {
	go heartbeat(ctx)
	sugar, err := sugarFromContext(ctx)
	if err != nil {
		return err
	}

	routerErr := initRouter(ctx)
	if routerErr != nil {
		sugar.Error("Failed to initialize router: %v", routerErr)
		return routerErr
	}
	return nil
}

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	ctx := context.Background()
	ctx, err := createAndSaveLoggers(ctx)
	if err != nil {
		os.Exit(1)
	}

	sugar, err := sugarFromContext(ctx)
	if err != nil {
		os.Exit(1)
	}

	sugar.Info("Starting application!!")

	runErr := run(ctx, os.Args)
	if runErr != nil {
		sugar.Error("%s\n", err)
		os.Exit(1)
	}
}
