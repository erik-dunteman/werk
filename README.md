# werk
A basic template for Golang Machinery distributed worker queue.

Werk is a minimal refactoring of [RichardKnop/machinery](https://github.com/RichardKnop/machinery), an "asynchronous task queue/job queue based on distributed message passing." For those familiar, it's the Go version of Python's [Celery](https://github.com/celery/celery).

Note: this is template, not a go package. Do not try to `go get` it.

#### What is an asynchronous task queue?

An asynchronous task queue uses a central message datastore (in this case, Redis) to allow one or more "callers" to register jobs in a central location, and one or more distributed "workers" to pull jobs from that queue, run them, and return the response.

#### Why use an asynchronous task queue?

Some http server jobs take longer than a single http request, so you need to instead start the job, run it in the background, and retrieve the output in the future. 

In horizintally scaled systems, there's no guarantee that the load balancer will hit the same server instance when the client attempts to retrieve the output of a particular job.

Using Redis to register jobs ensures that all http servers (the "callers") and the distributed "workers" can see all of the jobs, regardless of how many replicates are running.

#### Why this repo?
This refactor was done because the Machinery [examples](https://github.com/RichardKnop/machinery/tree/master/example/redis) launch both the worker and the caller processes from the same mixed codebase, making it difficult to view the two as separate, independently scaleable services. The callers and workers can and should be ran on multiple machines.

Werk splits the worker and caller logic into two separate folders, with a shared "machine" package that manages the shared components.

## How to run it

#### 1) Start the Redis server, which each worker and caller will connect to.
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
In this repo, it will place a new job into the queue, briefly wait for the worker to finish, and retrieve the output.
```bash
cd caller
go run caller.go
```

## How to customize it
