package main

import (
	"context"
	"time"

	tracers "github.com/RichardKnop/machinery/example/tracers"
	log "github.com/RichardKnop/machinery/v1/log"
	tasks "github.com/RichardKnop/machinery/v1/tasks"
	uuid "github.com/google/uuid"
	opentracing "github.com/opentracing/opentracing-go"
	opentracing_log "github.com/opentracing/opentracing-go/log"

	machine "github.com/erik-dunteman/werk/machine"
)

func main() {

	/*
	 * Fire up the caller's connection to redis
	 */

	cleanup, err := tracers.SetupTracer("sender")
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	// Start the server
	server, err := machine.StartServer()
	if err != nil {
		log.FATAL.Fatalln("Unable to start the front machine server:", err)
	}

	// Create a span and context to track through the entire worker call
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()
	batchID := uuid.New().String()
	span.SetBaggageItem("batch.id", batchID)
	span.LogFields(opentracing_log.String("batch.id", batchID))
	log.INFO.Println("Starting batch:", batchID)

	/*
	 * Now, we send our tasks!
	 */

	/*
	 * Task example 1: Sleep
	 */

	// First, init the task by task name and proper input argument
	thisTask := tasks.Signature{
		Name: "sleep", // must be same name as in machine.Tasks
		Args: []tasks.Arg{
			{
				Type:  "int",
				Value: 1,
			},
		},
	}

	// Send that task
	asyncResult, err := server.SendTaskWithContext(ctx, &thisTask)
	if err != nil {
		log.FATAL.Fatalln("Could not send task:", err.Error())
	}

	// Now normally here, we'd retrieve the result with a separate http server call,
	// but we'll simulate waiting a bit to retrieve it below

	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		log.FATAL.Fatalln("Getting task result failed with error:", err.Error())
	}
	log.INFO.Printf("Output = %v\n", tasks.HumanReadableResults(results))

	/*
	 * Task example 2: Hello
	 */

	// First, init the task by task name and proper input arguments
	thisTask = tasks.Signature{
		Name: "hello", // must be same name as in machine.Tasks
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: "cool user",
			},
		},
	}

	// Send that task
	asyncResult, err = server.SendTaskWithContext(ctx, &thisTask)
	if err != nil {
		log.FATAL.Fatalln("Could not send task:", err.Error())
	}

	// Now normally here, we'd retrieve the result with a separate http server call,
	// but we'll simulate waiting a bit to retrieve it below

	results, err = asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		log.FATAL.Fatalln("Getting task result failed with error:", err.Error())
	}
	log.INFO.Printf("Output = %v\n", tasks.HumanReadableResults(results))

}
