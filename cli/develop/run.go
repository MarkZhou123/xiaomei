package develop

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bughou-go/xiaomei/cli/setup/appserver"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Run() {
	if err := build(); err != nil {
		return
	}
	tail := tailLog()

	config.Log(`starting.`)
	if app, err := cmd.Start(cmd.O{}, config.App.Bin()); err != nil {
		return
	} else {
		appserver.WaitPort(os.Getpid(), app.Process.Pid)
		app.Wait()
	}

	tail.Process.Kill()
}

func tailLog() *exec.Cmd {
	appLog := filepath.Join(config.App.Root(), `log/app.log`)
	appErr := filepath.Join(config.App.Root(), `log/app.err`)
	cmd.Run(cmd.O{Panic: true}, `touch`, `-a`, appLog, appErr)
	tail, _ := cmd.Start(cmd.O{Panic: true}, `tail`, `-fqn0`, appLog, appErr)
	return tail
}

func build() error {
	config.Log(`building.`)
	if cmd.Ok(cmd.O{Env: []string{`GOBIN=` + config.App.Root()}}, `go`, `install`) {
		return nil
	}
	return errors.New(`build failed.`)
}
