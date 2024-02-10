# Gopts

Gopts is a small package to handle options in Go heavily inspired by Rust's
Option enum.

basically, it's a wrapper around

```go
if x != nil {
    return *x
}
```

## Import

```bash
go get github.com/karim-w/gopts
```

## Usage

```go
package main

import (
    "log"
    "github.com/karim-w/gopts"
)

func getOption() gopts.Option[int] {
    // Create an option with a value
    return gopts.Some(5)
}

func getNothing() gopts.Option[int] {
    // Create an option with no value
    return gopts.None[int]()
}

func main() {
    // an option
    val := getOption()

    // check if there is a value
    if val.IsSome() {
        // pull the value
        num := val.Unwrap()

        log.Println("Value is", num)
    }

    // an option with no value
    val = getNothing()

    // check if there is nothing
    if val.IsNone() {
        log.Println("No value")
    }


    defer func() {
        if r := recover(); r != nil {
            log.Println("Panic:", r)
        }
    }()


    // when pulling a value from an option with no value
    // it will panic
    val.Unwrap()

}
```

## License

BSD 3-Clause License

## Author

karim-w

## Contributing

Pull requests are welcome. For major changes, please open an issue first to
discuss what you would like to change.
