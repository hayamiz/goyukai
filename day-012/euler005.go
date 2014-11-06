//usr/bin/env go run $0 $@ ; exit

package main

import (
    "fmt"
    "math"
)

func factorize(x int) []int {
    ret := make([]int, 0, 100)

    if x == 1 {
        return []int{1}
    } else if x == 2 {
        return []int{2}
    } else if x == 3 {
        return []int{3}
    }

    for x % 2 == 0 {
        x /= 2
        ret = append(ret, 2)
    }

    for div := 3; div < int(math.Sqrt(float64(x))); div += 2 {
        for x % div == 0 {
            x /= div
            ret = append(ret, div)
        }
    }

    if x != 1 {
        ret = append(ret, x)
    }

    return ret
}

func main() {
    ret := 1

    for i := 1; i <= 20; i++ {
        factors := factorize(i)
        fmt.Println(i, factors)
        for _, fac := range factors {
            if ret % fac != 0 {
                ret *= fac
            }
        }
    }

    fmt.Println(ret)
}
