# Semaphores

Semaphores are a very general synchronization mechanism that can be used to implement mutexes, limit access to multiple resources, solve the readers-writers problem, etc.

There is no semaphore implementation in Go's sync package, but they can be emulated easily using buffered channels:

* the capacity of the buffered channel is the number of resources we wish to synchronize
* the length (number of elements currently stored) of the channel is the number of resources current being used
* the capacity minus the length of the channel is the number of free resources (the integer value of traditional semaphores)

We don't care about what is stored in the channel, only its length; therefore, we start by making a channel that has variable length but 0 size (in bytes):

```go
type empty {}
type semaphore chan empty
```
We then can initialize a semaphore with an integer value which encodes the number of available resources.  If we have N resources, we'd initialize the semaphore as follows:

```go
sem = make(semaphore, N)
```
Now our semaphore operations are straightforward:

```go
// acquire n resources
func (s semaphore) P(n int) {
    e := empty{}
    for i := 0; i < n; i++ {
        s <- e
    }
}

// release n resources
func (s semaphore) V(n int) {
    for i := 0; i < n; i++ {
        <-s
    }
}
```
This can be used to implement a mutex, among other things:

```go
/* mutexes */

func (s semaphore) Lock() {
    s.P(1)
}

func (s semaphore) Unlock() {
    s.V(1)
}

/* signal-wait */

func (s semaphore) Signal() {
    s.V(1)
}

func (s semaphore) Wait(n int) {
    s.P(n)
}
```
