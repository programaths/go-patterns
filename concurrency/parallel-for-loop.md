# Parallel For-Loop


Some languages provide a parallel for-loop (e.g. Sun's Fortress) which can simplify programming parallel algorithms.  Go doesn't support parallel for-loops as a separate construct, but they are easy to implement using goroutines.

## Usage


```go
type empty {}
...
data := make([]float32, N)
res := make([]float32, N)
var wg sync.WaitGroup
...
for i,xi := range data {
    wd.add(1)
    go func (i int, xi float32) {
        defer wg.Done()
        res[i] = doSomething(i,xi)        
    } (i, xi);
}
// wait for goroutines to finish
wg.Wait()
```

Notice the use of the anonymous closure.  The current i,xi are passed to the closure as parameters, masking the i,xi variables from the outer for-loop.  This allows each goroutine to have its own copy of i,xi; otherwise, the next iteration of the for-loop would update i,xi in all goroutines. On the other hand, the res array is not passed to the anonymous closure, since each goroutine does not need a separate copy of the array (or slice of the array).  The res array is part of the closure's environment but is not a parameter.

A somewhat practical example:
```go
func VectorScalarAdd (v []float, s float32) {
    var wg sync.WaitGroup
    for i,_ := range v {
        wg.Add(1)
        go func (i int) {
            defer wg.Done()
            v [i] += s
        } (i)
    }
    wg.Wait()
}
```
## For-Loops and Futures

When implementing a function which contains one big parallel for-loop (like the VectorScalarAdd example above), you can increase parallelism by returning a future rather than waiting for the loop to complete.

## Common Mistakes

It is easy to be overly dependent on channels in Go.  I've seen code like the following in several places:
```go
xi := make(float32 chan);
out := make(float32 chan);
// start N goroutines
for _,_ := range data {
    go func () {
        xi := <-xch;
        out <- doSomething(xi);
    }
}
// send input to each goroutine
for _,xi := range data {
    xch <- xi;
}
// collect results of each goroutine
for _,_ := range data {
    res := <-out;
    ....
}
```
In addition to being more verbose, it is inefficient because of the extra set-up and tear-down for-loops.  Notice too that this isn't very parallel: most of the time spent by each goroutine will be waiting for xch to be ready for reading or for res to be ready for writing.  This can be "solved" by using channels of capacity N to prevent blocking, or by creating separate channels for each goroutine:

```go
for _,xi := range data {
    xch := make(float32 chan)
    go func () {
        xi := <- xch;
        out <- doSomething(xi)
    }
    xch <- xi;
}
....
```
Of course, making N channels is much less efficient than passing parameters on the stack.  In summary, we need a channel for synchronization purposes (used as a semaphore) when implementing a parallel for-loop, but we do not need to communicate with goroutines through channels when the stack works perfectly well.
