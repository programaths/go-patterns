# Constructors
Go doesn't support constructors, but constructor-like factory functions are easy to implement:

```go
package matrix
function New(rows, cols int) *matrix {
    m := new(matrix)
    m.rows = rows
    m.cols = cols
    m.elems = make([]float, rows*cols)
    return m
}
```
To prevent users from instantiating uninitialized objects, the struct can be made private.

```go
package main
import "matrix"
wrong := new(matrix.matrix)    // will NOT compile (matrix is private)
right := matrix.New(2,3) // ONLY way to instantiate a matrix

```
## Initializers

Go aims to prevent unnecessary typing.  Oftentimes Go code is shorter and easier to read than object-oriented languages.  Compare the use of factory functions and initializers in Go to the use of constructors in Java:

```go
matrix := New(10, 10)
pair := &Pair{"one", 1}
```
```java
// Java:
Matrix matrix = new Matrix(10, 10);
Pair pair = new Pair ("one", 1);
```
Additionally, Go "constructors" can be written succinctly using initializers within a factory function:

```go
function New(rows, cols, int) *matrix {
    return &matrix{rows, cols, make([]float, rows*cols)}
}
```
