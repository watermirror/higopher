package termination

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Flag whether the process is asked to terminate.
var terminatingFlag bool = false

// During this period, program can handle system signals including SIGTERM,
// SIGINT and SIGQUIT.
var terminatingTimeoutMilliseconds int = 0

// If the program exits by terminating timeout, this code will be the exit code.
var exitCodeOnTimeout int8 = 0

// EnableGracefulTermination enables the feature of graceful termination.
// - The parameter timeoutMilliseconds is a during period, within which the
//   program can handle some terminating tasks.
// - If the program exits by terminating timeout, parameter timeoutExitCode will
//   be the exit code.
func EnableGracefulTermination(timeoutMilliseconds int, timeoutExitCode int8) {
	// Set timeout.
	if terminatingTimeoutMilliseconds < 0 {
		panic(1)
	}
	terminatingTimeoutMilliseconds = timeoutMilliseconds
	exitCodeOnTimeout = timeoutExitCode

	// Prepare system signal channel to listen signals sent to program.
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(
		signalChannel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Start to listen the system signals.
	go doListen(signalChannel)
}

// AllowTermination allows the program terminate by the sent system signal, if
// any.
// - If it was a terminating signal sent, the program will terminate with an
//   exit code equaling the parameter exitCode.
func AllowTermination(exitCode int) {
	if terminatingFlag == true {
		os.Exit(exitCode)
	}
}

// Listen the system signals.
func doListen(channel chan os.Signal) {
	// Pend gorouting to wait for signals.
	<-channel

	// Mark terminating flag.
	terminatingFlag = true

	// Wait for other part of the program to accomplish some terminationg
	// prosedures.
	time.Sleep(time.Duration(terminatingTimeoutMilliseconds) * time.Millisecond)
	os.Exit(int(exitCodeOnTimeout))
}
