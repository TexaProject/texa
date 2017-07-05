# Turing test EXtended Architecture (TEXA): A method for interacting, evaluating and assessing AI

TEXA is a novel framework backed by a mathematical oriented theory(Texa Theory, refer to /texa-docs) for interacting, evaluating and assessing AIs.


# LIBRARIES USED:

The system uses a number of open source projects to work properly:

* [texalib] - Special Math library written from scratch for TEXA
* [texajson] - Dedicated JSON interpreter library that serves as interface to  native Datastore.
* [JS] - evented I/O for the backend
* [jQuery] - duh
* [ElizaBOT-JS] - Javascript implementation of the ELIZA specification by Weizenbaum, 1966. Special thanks to Landsteiner!


# Installation:

TEXA requires [Go Lang](https://golang.org/)  v1.7+ to run.

Install the dependencies, devDependencies and start the server.

``$ go version``

``$ cd $GOPATH/src/``

``$ git clone https://github.com/TexaProject/texa.git``

``$ go get -u https://github.com/TexaProject/texajson.git``

``$ go get -u https://github.com/TexaProject/texalib.git``

``$ cd $GOPATH/src/texa/``

``$ go run main.go``


### TODO (Future Work)

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