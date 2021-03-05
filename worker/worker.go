package main

import (
	tracers "github.com/RichardKnop/machinery/example/tracers"
	log "github.com/RichardKnop/machinery/v1/log"
	tasks "github.com/RichardKnop/machinery/v1/tasks"

	machine "github.com/erik-dunteman/werk/machine"
)

func main() {
	consumerTag := "machinery_worker"

	cleanup, err := tracers.SetupTracer(consumerTag)
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	server, err := machine.StartServer()
	if err != nil {
		log.FATAL.Fatalln("Unable to start the worker server", err)
	}

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker(consumerTag, 0)

	// Here we inject some custom code for error handling,
	// start and end of task hooks, useful for metrics for example.
	errorhandler := func(err error) {
		log.ERROR.Println("I am an error handler:", err)
	}

	pretaskhandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am a start of task handler for:", signature.Name)
	}

	posttaskhandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am an end of task handler for:", signature.Name)
	}

	worker.SetPostTaskHandler(posttaskhandler)
	worker.SetErrorHandler(errorhandler)
	worker.SetPreTaskHandler(pretaskhandler)

	err = worker.Launch()
	if err != nil {
		log.FATAL.Fatalln("Unable to launch the worker", err)
	}
}
