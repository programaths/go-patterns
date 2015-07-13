# Generators

Generators are functions that return the next value in a sequence each time the
function is called:

```go
generateInteger() => 0
generateInteger() => 1
generateInteger() => 2
....
```

## Parallelism

If the generator's task is computationally expensive, the generator pattern can allow a consumer to run in parallel with a generator as it produces the next value to consume.  For example, the goroutine behind the "produce" generator below can execute in parallel with "consume".

```go
for {
    consume(produce())
}
```
In some cases, the generator itself can be parallelized.  When a generator's task is computationally expensive and can be generated in any order (constrast with iterators), then the generator can be parallelized internally:

```go
func generateRandomNumbers (n int) {
    ch := make (chan float)
    sem := make (semaphore, n)

    for i := 0; i < n; i++ {
        go func () {
            ch <- rand.Float()
            close(ch)
        } ()
    }

    // launch extra goroutine to eventually close ch
    go func () {
        sem.Wait(n)
        close(ch)
    }

    return ch
}

```

## Usage

The following for-loop will print 100 random numbers.  The random numbers are generated in parallel and arrive in a random order.  Since the order doesn't matter, this prevents us from blocking unnecessarily.

```go
for x := range generateRandomNumbers(100) {
    fmt.Println(x)
}
```
Note that generating random numbers in this way isn't very practical, since the parallelizing random number generator probably isn't worth the overhead of spawning so many goroutines.
