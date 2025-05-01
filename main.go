package main

import (
	"awesomeProject/utils"
	"context"
	"os"
)

var (
	createAndSaveLoggersFunc = utils.CreateAndSaveLoggers
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func run(ctx context.Context, args []string) error {
	return nil
}

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	ctx := context.Background()
	ctx, err := createAndSaveLoggersFunc(ctx)
	if err != nil {
		os.Exit(1)
	}

	sugar, err := utils.SugarFromContext(ctx)
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
