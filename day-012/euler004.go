//usr/bin/env go run $0 $@ ; exit

package main

import (
    "fmt"
)

func isPalindrome(x int) bool {
    x_str := fmt.Sprintf("%d", x)
    for i := 0; i < len(x_str) / 2; i++ {
        if x_str[i] != x_str[len(x_str) - 1 - i] {
            return false
        }
    }
    return true
}

func main() {
    max_val := 0

    for x := 100; x < 1000; x++ {
        for y := 100; y < 1000; y++ {
            if isPalindrome(x * y) {
                if x * y > max_val {
                    max_val = x * y
                }
            }
        }
    }

    fmt.Println(max_val)
}
