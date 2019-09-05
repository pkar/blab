# blab

Testing out a tcp chat server

Run the server

```
$ go run cmd/blab/main.go -h
  -host string
    	The host interface to listen on (default "localhost")
  -logs string
    	The directory to write chat logs to
  -port int
    	The port the server accepts connections for clients (default 7777)


$ go run cmd/blab/main.go -logs="logs"
2019/09/05 09:57:03 server.go:40: INFO: listening localhost 7777
```

Connect a TELNET client

```
$ telnet localhost 7777
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.

----------------------
\help: print help
\join <roomname>: enter the name of the room to join, or create a new one
\list: list all available rooms
\name <name>: change the user name
\quit: quit
----------------------
$ Enter name: bob
$ Hi bob, join a room first with command `\join`, then enter text
hi
$ Join a room to send a message
\join test
2019-09-05T09:57:17-07:00] [test] bob joined...
hello
2019-09-05T09:57:21-07:00] [test] (bob) hello
2019-09-05T09:57:30-07:00] [test] bill joined...
2019-09-05T09:57:38-07:00] [test] (bill) hello
2019-09-05T09:57:46-07:00] [test] bill has left..
\quit
2019-09-05T09:57:54-07:00] [test] bob has left..
Connection closed by foreign host.
```
