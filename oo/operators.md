# Operators

An operator is a unary or binary function which returns a new object and does not modify its parameters. In C++, special infix operators (+, -, *, etc) can be overloaded to support math-like syntax. Since Go does not support operator overloading, ordinary functions must be used to manipulate user types. There are two major approaches to dealing with this limitation: operators as functions, and operators as methods.

## Operators as Functions

Use package-level functions to operate on one or two parameters and return a new object.  This supports the following idioms:

```go
c := complex.Add(a, b)
m := matrix.Add(m1, matrix.Mult(m2, m3))

```

These functions should be implemented within a package related to the objects on which they operate.

Go does not support function overloading, so oftentimes several variations of a operator may be required:

```go
func addSparseToDense (a *sparseMatrix, b *denseMatrix) *denseMatrix
func addDenseToDense (a *denseMatrix, b *denseMatrix) *denseMatrix
func addSparseToSparse (a *sparseMatrix, b *sparseMatrix) *sparseMatrix
```
This is a rather clumsy API, so make these private and provide a single public function with a type-switch:

```go
func Add (a Matrix, b Matrix) Matrix {
    switch a.(type) {
    case sparseMatrix:
        switch b.(type) {
        case sparseMatrix:
            return addSparseToSparse(a.(sparseMatrix), b.(sparseMatrix))
        case denseMatrix:
            return addSparseToDense(a.(sparseMatrix), b.(denseMatrix))
    ....
    default:
        // unsupported arguments
        ....
    }
}
```
Here a Matrix interface is made public and the Add function can operate on any combination of supported parameters.

## Operators as Methods

Another approach is to use method chaining:

```go
m := m2.Times(m3).Plus(m1)
```
Each method returns a new object which becomes the receiver of the next method call.  Implementation is similar to the previous approach but involves a receiver:
```go
func (a *denseMatrix) Plus(b Matrix) Matrix
func (a *sparseMatrix) Plus(b Matrix) Matrix
```
Since Go allows methods to be "overloaded" (for lack of a better term) based on the receiver, it is possible to provide separate implementations for each type of receiver.  The correct implementation will be selected at runtime based on V-tables.  Each implementation can perform a type-switch if necessary:
```go
func (a *denseMatrix) Plus(b Matrix) Matrix {
    switch b.(type) {
    case sparseMatrix:
....
```
##Best Practices

Use an interface and polymorphism when you want to use operators:

```go
type Algebraic interface {
    Plus(b Algebraic) Algebraic;
    Minus(b Algebraic) Algebraic;
    Times(b Algebraic) Algebraic;
...
}
func (a *complex) Plus(b Algebraic) Algebraic
func (a *rational) Minus(b Algebraic) Algebraic
```
Each type which implements the Algebraic interface above will allow for method chaining.  Each method implementation should use a type-switch to provide optimized implementations based on the parameter type.  Additionally, a default case should be specified which relies only on the methods in the interface:

```go
func (a *denseMatrix) Plus(b Algebraic) Algebraic {
    switch b.(type) {
    case sparseMatrix:
        return addDenseToSparse(a, b.(sparseMatrix))
    default:
        for x in range b.Elements() ....
....
```
If a generic implementation cannot be implemented using only the methods in the interface, you probably are dealing with classes that are not in the same "algebra", and this operator pattern should be abandoned.  For example, it does not make sense to write a.Plus(b) if 'a' is a set and 'b' is a matrix; therefore, it will be difficult to implement a generic a.Plus(b) in terms of set and matrix operators.  In this case, split your package in two and provide separate AlgebraicSet and AlgebraicMatrix interfaces.
