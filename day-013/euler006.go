//usr/bin/env go run $0 $@ ; exit

package main

import (
    "fmt"
)

func main() {
    sqsum := 0
    sum := 0

    for i := 1; i <= 100; i++ {
        sqsum += i * i
        sum += i
    }

    fmt.Println(sum * sum - sqsum)
}
