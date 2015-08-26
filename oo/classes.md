# Classes

Go doesn't require an explicit class definition as Java, C++, C#, etc do.  Instead, a "class" is implicitly defined by providing a set of "methods" which operate on a common type.  The type may be a struct or any other user-defined type.  For example:

```go
type Integer int

func (i *Integer) String() string {
    return strconv.itoa(i)
}
```
... is analogous to:

```
class Integer {
    public int i;
    public String toString() { return Integer.toString(i); }
}
```
... in Java.
