# Turing test EXtended Architecture (TEXA): A method for interacting, evaluating and assessing AI

TEXA is a novel framework backed by a mathematical oriented theory(Texa Theory, refer to /texa-docs) for interacting, evaluating and assessing AIs.


## LIBRARIES USED:

The system uses a number of open source projects to work properly:

* [texalib] - Special Math library written from scratch for TEXA
* [texajson] - Dedicated JSON interpreter library that serves as interface to  native Datastore.
* [JS] - evented I/O for the backend
* [jQuery] - duh
* [ElizaBOT-JS] - Javascript implementation of the ELIZA specification by Weizenbaum, 1966. Special thanks to Landsteiner!


## Installation:

### Go:

TEXA is written in [Go](http://golang.org) v1.7+. If you don't have a Go
development environment, please [set one up](http://golang.org/doc/code.html).

### Redis: 

Texa using the [texalib]
to store the slabs, CatPages and MtsPages, in [redis](https://redis.io/)
[Redis quick start](https://redis.io/topics/quickstart#redis-quick-start).

### MongoDB:

TEXA using [MongoDB](https://docs.mongodb.com/manual/installation) to persist the conversation between 
AI and human for every session.

### Build on a local OS:

Make sure your go version is 1.7+

``$ go version``

Move to your 

``$ cd $GOPATH/src/``

Clone the texa project from github

``$ git clone https://github.com/TexaProject/texa.git``

Navigate to the texa project

``$ cd $GOPATH/src/texa/``

Run make file, by the following order

```sh
   make install
   make test
   make run

   make fmt
   #before you push your changes, to format the code base
```

follow the below steps, to set up your local git config to contribute

```sh
git remote add upstream https://github.com/TexaProject/texa.git
# or: git remote add upstream git@github.com:TexaProject/texa.git

# Never push to upstream master
git remote set-url --push upstream no_push

# Confirm that your remotes make sense:
git remote -v
```

### Branch

Get your local master up to date:

```sh
git fetch upstream
git checkout master
git rebase upstream/master
```

Branch from it:
```sh
git checkout -b mytexa
```

Then edit code on the `mytexa` branch.

## Unit test and bench mark test:

Run unit test on package level

``go test texa/storage``

Run Benchmark test on package level

``go test go test texa/storage -bench=.``

Run a unit test, navigate to package 

``go test -run=TestAddtoMongo``

Run a benchmark test

``go test -run=BenchmarkAddtoMongo -bench=.``

## Profiling Texa App:

Profiling Texa App by using Go [pprof package](https://golang.org/pkg/net/http/pprof/)
To use pprof, link this package into app:

``import _ "net/http/pprof"``

### Pprof:

[pprof](https://github.com/google/pprof) is a tool for visualization and analysis of profiling data.
Install the belos dependency to build the pprof tool

```go get -u github.com/google/pprof```

### Pprof sample:

#### Profiling http request:

Testing the sample load by manual script file.

```sh
   for i in {1..1000}; do curl -X POST http://localhost:3030/texa -d IntName=pan -d scoreArray=[0,1,1] 
   -d SlabName=sports -d slabSequence=[sports,sports,sports] 
   -d chatHistory=[hi,hello! how can I help you,how is weather today, do you wanna know more,yes,its nice talking to you] -d timeStamp=32465466754; done
```

Befor run the above script start to monitor the profiling of Texa app by 

```go tool pprof http://localhost:3030/debug/pprof/profile?seconds=30```

Profiling creates a pprof file eg : `pprof.main.samples.cpu.01.pb.gz`
Visualize and analysis by pprof tool 

```pprof -http=:6060 pprof.main.samples.cpu.01.pb.gz```

It starts pprof tool app on port 6060 and visualizing the profiling data including flame graph

#### Profiling benchmark test:

Below command will execute the benchmark test, and it will create a binary `storage.test`
and creates `profile.out` file.

```go test -run=BenchmarkAddtoMongo -bench=.```

Then, visualize and analysis the profiling data by pprof tool, pprof starts a service in 8081. 

```pprof -http=localhost:8081 storage.test profile.out```


## TODO (Future Work)

- Stability offered only for Unit Test Instance. Can support multiple cases without SLA.
- Lacks complete support for non-Eliza AIs(non-JS REF).
- APIs can be exposed to build use-cases such as ranking apps etc.
- Feel free to try new ideas!


License
----

Apache 2.0 on the demonstrated work.
Derived work carry respective Licenses. Please refer the links.

[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen. Thanks SO - http://stackoverflow.com/questions/4823468/store-comments-in-markdown-syntax)

   [texalib]: <https://github.com/TexaProject/texalib>
   [texajson]: <https://github.com/TexaProject/texajson>
   [JS]: <http://nodejs.org>
   [Twitter Bootstrap]: <http://twitter.github.com/bootstrap/>
   [jQuery]: <http://jquery.com>
   [ElizaBOT-JS]: <http://www.masswerk.at/elizabot/>