# Futures

Sometimes you know you need to compute a value before you need to actually use the value.  In this case, you can potentially start computing the value on another processor and have it ready when you need it.  This is the idea behind futures.

Futures are easy to implement via closures and goroutines.  The idea is similar to generators, except a future needs only to return one value.

## Usage

```go
func InverseProduct (a Matrix, b Matrix) {
    a_inv := Inverse(a);
    b_inv := Inverse(b);
    return Product(a_inv, b_inv);
}
```
In the above contrived example, it is known initially that the inverse of both 'a' and 'b' must be computed.  Why should the program wait for a_inv to be computed before starting b_inv?  These Inverse computations can be done in parallel.  On the other hand, the call to Product needs to wait for both a_inv and b_inv to finish.  This can be implemented as follows:

```go
func InverseProduct (a Matrix, b Matrix) {
    a_inv_future := InverseFuture(a);
    b_inv_future := InverseFuture(b);
    a_inv := <-a_inv_future;
    b_inv := <-b_inv_future;
    return Product(a_inv, b_inv);
}
```
In this improved version, the InverseFuture function launches a goroutine to perform the inverse computation and immediately returns a channel which will eventually hold the future value:
```go
func InverseFuture (a Matrix) {
    future := make (chan Matrix);
    go func () { future <- Inverse(a)  }();
    return future;
}
```
The Inverse is computed asynchronously and potentially in parallel.

## Futures in APIs

When developing a computationally intensive package, it may make sense to design the entire API around futures.  The futures can be used within your package while maintaining a friendly API.  In addition, the futures can be exposed through an asynchronous version of the API.  This way the parallelism in your package can be lifted into the user's code with minimal effort:

```go
package "matrix"
// futures used internally
type futureMatrix chan Matrix;

// API remains the same
func Inverse (a Matrix) Matrix {
    return <-InverseAsync(promise(a))
}

func Product (a Matrix, b Matrix) Matrix {
    return <-ProductAsync(promise(a), promise(b))
}

// expose async version of the API
func InverseAsync (a futureMatrix) futureMatrix {
    c := make (futureMatrix);
    go func () { c <- inverse(<-a) } ();
    return c
}

func ProductAsync (a futureMatrix) futureMatrix {
    c := make (futureMatrix);
    go func () { c <- product(<-a) } ();
    return c
}

// actual implementation is the same as before
func product (a Matrix, b Matrix) Matrix {
    ....
}

func inverse (a Matrix) Matrix {
    ....
}

// utility fxn: create a futureMatrix from a given matrix
func promise (a Matrix) futureMatrix {
    future := make (futureMatrix, 1);
    future <- a;
    return future;
}
```

The above package can be used just as before:

```go
package "main"
func InverseProduct (a Matrix, b Matrix) {
    a_inv := Inverse(a);
    b_inv := Inverse(b);
    return Product(a_inv, b_inv);
}
```
...or asynchronously:
```go
package "main"
func InverseProduct (a Matrix, b Matrix) {
    a_inv_future := InverseAsync(a);
    b_inv_future := InverseAsync(b);
    a_inv := <-a_inv_future;
    b_inv := <-b_inv_future;
    return Product(a_inv, b_inv);
}
```
Either way, we've added more parallelism behind the scenes without rewriting the underlying algorithms.
