# Singleton

The singleton pattern is used to ensure that there is only one instance of a resource shared by many consumers. For more information on the pattern, take a look at the [singleton](https://en.wikipedia.org/wiki/Singleton_pattern) pattern page on Wikipedia.

This example is adapted from this [blog post by Marcio Castilho](http://marcio.io/2015/07/singleton-pattern-in-go/)

## Usage

The typical implementation of a singleton usually involves the use of a `GetInstance()` method. Invoking this method will first check to see if an instance of the resource has already been created. If it has, that instance will be returned. If not, a new one will be created and returned. 

The hard part here is to ensure that your GetInstance() method is thread-safe. The approach outlined by Marcio uses the Go `sync.Once`  type to achieve a thread-safe implementation with very low overhead:

```go
package singleton

import (
    "sync"
)

type singleton struct {
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
    once.Do(func() {
        instance = &singleton{}
    })
    return instance
}
```

Here the `once.Do()` method wraps an anonymous function in which the singleton instance is created. Because `sync.Once` guarantees that the `Do()` method will only ever be called one time your instance value will never be allocated twice.
