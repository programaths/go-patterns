# Iterators

A clever and convenient use of Go's channels is for the construction of iterators.

## Usage

Iterators in Go are supported by a natural syntax:

```go
for x := range channel {
    fmt.Println(x)
}
```
The `for...range` construct reads from the given channel until the channel is closed.  Obviously, another goroutine must be writing to the channel and must close the channel when it is done writing.  This goroutine is called a generator in the general case and specifically an "iterator" when it produces elements of a container.

Iterating over a container looks like:

```go
for x := range container.Iter() { ...
```
## Implementation

The Iter() method of a container returns a channel for the calling for-loop to read from.  A typical Iter() implementation looks like:

```go
func (c *container) Iter () <-chan item {
    ch := make(chan item);
    go func () {
        for i := 0; i < c.size; i++ {
            ch <- c.items[i]
        }
    } ();
    return ch
}
```
Inside the goroutine, a for-loop iterates over the elements in the container.  For tree or graph algorithms, this simple for-loop could be replaced with a depth-first search, for example.

## Lack of Parallelism

While the above iterator employs a channel and two goroutines (which may run in separate threads), the relationship between the two goroutines is such that one is usually blocking on the other.  Only one processor will be employed most of the time.  This problem can be ameliorated by using a channel with a buffer size greater than 0.  For example, with a buffer of size 100, the iterator can produce at least 100 items from the container before blocking.  If the consumer goroutine is running on a separate processor, it is possible that neither goroutine will ever block.

Since the number of items in the container is generally known, it makes sense to use a channel with enough capacity to hold all the items.  This way, the iterator will never block (though the consumer goroutine still might).  However, this effectively doubles the amount of memory required to iterate over any given container, so channel capacity should be limited to some maximum number.

## Efficiency Considerations

The process of copying each item from a container to and from a channel can make these sorts of iterators slower than other methods.  The speed of channels continues to improve, however.
