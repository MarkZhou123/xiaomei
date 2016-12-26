package config

import (
	"path/filepath"

	"github.com/bughou-go/xiaomei/utils/cmd"
)

var Deploy DeployConf

type DeployConf struct {
	conf deployConf
}

type deployConf struct {
	User      string `yaml:"user"`
	Root      string `yaml:"root"`
	GitAddr   string `yaml:"gitAddr"`
	GitBranch string `yaml:"gitBranch"`
}

func (d *DeployConf) Name() string {
	return App.Name() + `_` + App.Env()
}
func (d *DeployConf) Root() string {
	return d.conf.Root
}
func (d *DeployConf) Path() string {
	return filepath.Join(d.Root(), d.Name())
}
func (d *DeployConf) User() string {
	return d.conf.User
}

func (d *DeployConf) GitAddr() string {
	return d.conf.GitAddr
}

func (d *DeployConf) GitBranch() string {
	if d.conf.GitBranch != `` {
		return d.conf.GitBranch
	}
	d.conf.GitBranch, _ = cmd.Run(cmd.O{Output: true, Panic: true},
		`git`, `rev-parse`, `--abbrev-ref`, `HEAD`)
	return d.conf.GitBranch
}
