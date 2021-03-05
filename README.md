# werk
A basic template for Golang Machinery distributed worker queue.

Werk is a minimal refactoring of [RichardKnop/machinery](https://github.com/RichardKnop/machinery), an "asynchronous task queue/job queue based on distributed message passing." For those familiar, it's the Go version of Python's [Celery](https://github.com/celery/celery).

Note: this is template, not a go package. Do not try to `go get` it.

### What is an asynchronous task queue?

An asynchronous task queue uses a central message datastore (in this case, [Redis](https://redis.io/)) to allow one or more "callers" to register tasks in a central location, and one or more distributed "workers" to pull tasks from that queue, run them, and return the response.

### Why use an asynchronous task queue?

Some http server tasks take longer than a single http request, so you need to instead start the task, run it in the background, and retrieve the output in the future. 

In horizontally scaled systems, there's no guarantee that the load balancer will hit the same server instance when the client attempts to retrieve the output of a particular task. This is bad.

Using Redis to register tasks ensures that all http servers (the "callers") and the distributed "workers" can see all of the tasks, regardless of how many replicates are running.

### Why this repo?

This refactor was done because the Machinery [examples](https://github.com/RichardKnop/machinery/tree/master/example/redis) launch both the worker and the caller processes from the same mixed codebase, making it difficult to view the two as separate, independently scaleable services. The callers and workers can and should be ran on multiple machines.

Werk splits the worker and caller logic into two separate folders, with a shared "machine" package that manages the shared components.

## How to run Werk

#### 1) Clone this repo

Note: 
You will need to update the imports if you're running in a custom module.

The cloned repo will come with a go.mod file that registers this repo as `github.com/erik-dunteman/werk` by default.
Both `caller.go` and `worker.go` import `github.com/erik-dunteman/werk/machine` using this module name.

If you plan to rename your go module, change the `github.com/erik-dunteman/werk/machine` imports to your new name.

#### 2) Start the Redis server, which each worker and caller will connect to.

Via docker:
```bash
docker run -p 6379:6379 redis
```
If not using docker, follow [these instructions](https://redis.io/topics/quickstart) to start a redis server

#### 2) Start one (or more) workers

The `worker` folder contains the code to start a worker. Run:
```bash
cd worker
go run worker.go
```

#### 3) Start one (or more) callers

The `caller` folder contains the code to start the caller. 
In this repo, it will place a new task into the queue, briefly wait for the worker to finish, and retrieve the output.
```bash
cd caller
go run caller.go
```

## How to add tasks to Werk

#### 1) Edit the shared `machine` package

You define your task handlers in `tasks.go`. 
First, define the function as xyzHandler:
```go
// Do xyz stuff ...
func xyzHandler(t int) (bool, error) {
	// add logic here
	return true, nil
}
```
Second, register that task for caller/worker use by adding it to the Tasks map:
```go
// Define the tasks for external use
var Tasks = map[string]interface{}{
	"sleep": sleepHandler,
	"panic": panicHandler,
	"xyz": xyzHandler, // remember the name "xyz", that's how we call this task later.
}
```

This is the core of how Werk works. Both the callers and the workers will be importing this definition by accessing the Tasks map. When rolling changes to prod, the callers and workers must both be updated so that they share the same defined Tasks.

You can optionally edit server behavior (timeouts, redis endpoint, etc) in `config.go`.

#### 2) Edit the caller in `caller.go`

Call a newly defined task
```go
// xyz
thisTask := machine.InitTask()
thisTask := tasks.Signature{
	Name: "xyz", // must be same task name as in machine.Tasks
	Args: []tasks.Arg{ // send in the task handler arguments here. In this case, it's just one int.
		{
			Type:  "int",
			Value: 1,
		},
	},
}

asyncResult, err := server.SendTaskWithContext(ctx, &thisTask)
if err != nil {
}
  
// Now normally here, we'd retrieve the result with a separate http server call,
// but we'll simulate that by waiting shortly, then retrieving the result below.
  
results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
if err != nil {
}
```

## Doing more with Machinery
Consult with the original [Machinery](https://github.com/RichardKnop/machinery/) repo, and it's associated [example](https://github.com/RichardKnop/machinery/tree/master/example/redis) to extend this template.
