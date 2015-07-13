# Long-Poll Server

The term "comet" refers to web techniques that allow servers to push data to a web browser.  This is in contrast to HTTP's poll-only REST architecture.

One way to simulate push behavior via HTTP is to use a long poll: an HTTP request which waits for some even before responding.

Long-polling via XMLHttpRequest is straightforward.  Basically, the client polls the server like normal.  The server will not reply until a corresponding event has occurred.  When the server finally responds, a client-side callback is triggered, and another request is immediately.  This simple helper function demonstrates the idea:

```go
function longpoll(url, callback) {

  var req = new XMLHttpRequest ()
  req.open ('GET', url, true)

  req.onreadystatechange = function (aEvt) {
    if (req.readyState == 4) {
      if (req.status == 200) {
        callback(req.responseText)
        longpoll(url, callback)
      } else {
        alert ("long-poll connection lost")
      }
    }
  };

  req.send(null);
}
```
The server must be designed to wait until some event occurs before responding to the HTTP request.  This can be implemented via semaphores, locks, etc in other languages.  In Go, channels work great:
```go
func PollResponse(w http.ResponseWriter, req *http.Request) {
  io.WriteString(w, <-messages)
}
```

The HTTP response handler will not write the response until the message channel is non-empty.  This way, clients can long-poll for the message and be notified instantly when it arrives.

The message channel can be filled by other clients:
```go
func PushHandler(w http.ResponseWriter, req *http.Request) {

  body, err := ioutil.ReadAll(req.Body)

  if err != nil {
    w.WriteHeader(400)
  }

  messages <- string(body)
}
```
In the above response handler, the client is expected to POST a message, which will be instantly delivered to another long-polling client.

With named channels, it is possible to notify specific clients:

```go
var messages map[string] chan string = make(map[string] chan string)

func PushHandler(w http.ResponseWriter, req *http.Request) {

  rcpt := req.FormValue("rcpt")
  body, err := ioutil.ReadAll(req.Body)

  // check for bad requests
  if err != nil || rcpt == "" || req.Method != "POST" {
    w.WriteHeader(400)
  return
}

 ch := messages[rcpt]

 // new client?
 if ch == nil {
   ch = make (chan string)
   messages[rcpt] = ch
 }

 // store message
 if !(ch <- string(body)) {
   // channel full, or no one listening
   w.WriteHeader(503)
   return
 }
}
```

With non-zero-length channels, you can enable your server to queue messages for off-line clients.

However, there is a major problem with the above code.  The PollResponse handler might continue to block even after the client has disconnected.  Unless a message arrives for the client at some point, the response handler might continue to block indefinitely.

To solve this problem, you can launch an extra goroutine to send a timeout message to the response handler:

```go
func PollResponse(w http.ResponseWriter, req *http.Request) {

  timeout := make (chan bool)

  go func () {
    time.Sleep(30e9)
    timeout <- true
  }

  select {
    case msg := <-messages:
    io.WriteString(w, msg)
    case stop := <-timeout:
    return
  }
}
```
Perhaps when the http package matures, there will be a per-request channel for detecting "client disconnected" events.
