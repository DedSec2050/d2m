//go:build debug
// +build debug

package main

import (
	"os"

	"github.com/Techzy-Programmer/d2m/config/flow"
	"github.com/Techzy-Programmer/d2m/config/paint"
	"github.com/Techzy-Programmer/d2m/config/univ"
)

func init() {
	isDebug = true
}

func startDebug() {
	paint.Warn("Running in debug mode")

	flow.StartDeployment(&univ.DeploymentRequest{
		Branch: 				 "main",
		AutoSetupDeps:   true,
		LocalUser:       "risha",
		LocalParentPath: "Documents\\my-deploy",
		RepoPath:        "Node-Bug-Hunter/Hunter-UI",
		PreDeployCmds:   []string{"mkdir hello-world", "touch hello-world/index.html"},
		PostDeployCmds:  []string{"touch hello-world/deploy-success.txt"},
		FailOnCmdError:  false,
	})

	os.Exit(0)
}
