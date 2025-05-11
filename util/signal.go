package util

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Package util provides utility functions for handling signals and logging.
// It includes a function to wait for termination signals and execute a cleanup function.

// Fini is a function type that represents a cleanup function to be executed on termination.
type Fini func()

// WaitingToBeTerminate listens for termination signals (SIGINT, SIGTERM)
// and executes the provided cleanup function (fini) before exiting the program.
// It logs the received signal and the exit message using the provided logger (ilog).
// The function blocks until a signal is received, ensuring that the program can
// gracefully handle termination requests and perform any necessary cleanup before exiting.
func WaitingToBeTerminate(fini Fini, ilog *log.Logger) {
	// Create a channel to receive signals.
	sigChan := make(chan os.Signal, 1)

	// Notify the channel of specific signals.
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ilog.Println("Program is running. Press Ctrl+C to interrupt.")

	// Block until a signal is received.
	sig := <-sigChan

	// Handle the signal.
	switch sig {
	case syscall.SIGINT:
		ilog.Println("SIGINT received. Exiting...")
		fini()
		os.Exit(0)
	case syscall.SIGTERM:
		ilog.Println("SIGTERM received. Exiting...")
		fini()
		os.Exit(0)
	default:
		ilog.Println("Unexpected signal received.")
	}
}
