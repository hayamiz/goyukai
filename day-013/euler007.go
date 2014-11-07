//usr/bin/env go run $0 $@ ; exit

package main

import (
    "fmt"
    "math"
)

var prime_table []bool
var table_size int = 1000

func isPrime(n int) bool {
    if n < table_size {
        return prime_table[n]
    }

    search_limit := int(math.Sqrt(float64(n)))
    if n <= search_limit {
        panic(n)
    }

    for i := 0; i < table_size; i++ {
        if prime_table[i] == false {
            continue
        }
        if n % i == 0 {
            return false
        }
    }

    return true
}

func main() {
    prime_table = make([]bool, table_size)

    // build sieve of Eratosthenes
    prime_table[0] = false
    prime_table[1] = false
    prime_table[2] = true
    prime_table[3] = true
    for n := 5; n < table_size; n++ {
        prime_table[n] = true
    }
    for n := 3; n < table_size; n++ {
        if n & 1 == 0 {
            prime_table[n] = false
        } else {
            if prime_table[n] == true {
                for m := 2; m * n < table_size; m++ {
                    prime_table[m * n] = false
                }
            }
        }
    }

    prime_cnt := 0
    for n := 1; true ; n++ {
        if isPrime(n) {
            prime_cnt ++
            if prime_cnt == 10001 {
                fmt.Println(n)
                break
            }
        }
    }
}
