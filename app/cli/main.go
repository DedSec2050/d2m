package main

import (
	"flag"
	"log"
	"net"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/Techzy-Programmer/d2m/app/daemon"
	"github.com/Techzy-Programmer/d2m/cmd"
	"github.com/Techzy-Programmer/d2m/config/db"
	"github.com/Techzy-Programmer/d2m/config/univ"
	"github.com/Techzy-Programmer/d2m/internal/ipc"
	"github.com/urfave/cli/v2"
)

var Release string

func init() {
	univ.IsProd = (Release == "prod")
}

func main() {
	if !univ.IsProd {
		startDebug()
	}

	daemonFlag := flag.Bool("daemon", false, "Run as daemon")
	flag.Parse()

	if *daemonFlag {
		daemon.LaunchDaemon()
		return
	}

	// Ensure daemon is running before handling CLI commands
	pid := db.GetConfig[float64]("daemon.PID")

	if !isProcessRunning(pid) {
		ensureDaemonRunning()

		if !connectToDaemon() {
			panic("Unable to connect with daemon over TCP")
		}
	}

	app := &cli.App{
		Name:   "d2m",
		Usage:  "Managr your deployments with ease",
		Action: cmd.HandleInitCMD,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)
}

func isProcessRunning(pid float64) bool {
	process, err := os.FindProcess(int(pid))

	if err != nil {
		return false
	}

	if runtime.GOOS == "windows" {
		return checkTCPAlive()
	}

	sigErr := process.Signal(syscall.Signal(0))
	if sigErr == nil {
		return checkTCPAlive()
	}

	return false
}

// Function to check if it's really our own daemon service that's running with the given PID
func checkTCPAlive() bool {
	if !connectToDaemon() {
		return false
	}

	select {
	case <-univ.AliveChannel:
		return true

	case <-time.After(5 * time.Second):
		univ.CLIConn.Close()
		return false
	}
}

func connectToDaemon() bool {
	port := db.GetConfig[string]("daemon.Port")
	if port == "" {
		return false
	}

	conn, err := net.Dial("tcp", ":"+port)
	if err != nil {
		return false
	}

	go ipc.HandleConnection(conn)
	univ.CLIConn = conn
	return true
}
