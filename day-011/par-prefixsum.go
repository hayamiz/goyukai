//usr/bin/env go run $0 $@ ; exit

package main

import (
    "fmt"
    "runtime"
    "math/rand"
    "time"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    N := 1000000
    TRY := 5

    seq_time := 0.0
    par_time := 0.0

    for try := 0; try < TRY; try ++ {
        input := make([]int,  N)
        output := make([]int, len(input))
        par_output := make([]int, len(input))
        for i, _ := range input {
            input[i] = int(rand.Int31())
        }

        seq_time += exec_time(func() {
            prefixsum(input, output)
        })

        par_time += exec_time(func() {
            par_prefixsum(input, par_output)
        })

        // cross check
        for j, _ := range output {
            if output[j] != par_output[j] {
                panic(fmt.Sprintf("prefix sum mismatch: %d", j))
            }
        }
    }

    seq_time /= float64(TRY)
    par_time /= float64(TRY)

    fmt.Printf("#seq_time	par_time\n")
    fmt.Printf("%f	%f\n", seq_time, par_time)
}

func exec_time(body_proc func()) float64 {
    t0 := time.Now()
    body_proc()
    dt := time.Since(t0)
    return dt.Seconds()
}

func prefixsum(input []int, output []int) {
    for i, x := range input {
        if i > 0 {
            output[i] = output[i - 1] + x
        } else {
            output[i] += x
        }
    }
}

func par_prefixsum(input []int, output[]int) {
    // fmt.Println("par_prefixsum:", input)
    if len(input) == 1 {
        // fmt.Println("1 length input:")
        output[0] = input[0]
    } else {
        y := make([]int, len(input) / 2)
        z := make([]int, len(input) / 2)
        // fmt.Println("recur prefix sum: len(y) = ", len(y))
        parallel_for(0, len(y), func(i int){
            y[i] = input[2 * i] + input[2 * i + 1]
        })
        // fmt.Println("y: ", y, ", z:", z)
        par_prefixsum(y, z)
        parallel_for(0, len(input), func(i int){
            if i == 0 {
                output[0] = input[0]
            } else if i % 2 == 1 {
                output[i] = z[i / 2]
            } else {
                output[i] = z[(i - 1) / 2] + input[i]
            }
        })
    }
}


// loop for i := range [int]{i_low, i_low + 1, ... , i_up - 1}
// loop over [i_low, i_up)
func parallel_for(i_low int, i_up int, for_proc func(int)) {
    // fmt.Printf("parallel_for: i_low = %d, i_up = %d\n", i_low, i_up)
    ch := make(chan int, i_up - i_low)

    N := i_up - i_low
    partNum := runtime.NumCPU() * 2
    partLength := N / partNum
    if partLength < 4 {
        partLength = 4
    }

    ch = make(chan int, partNum)
    for idx := 0; idx < N; idx += partLength {
        // fmt.Printf("  parallel_for: dispatch %d ... %d\n", idx, idx + partLength)
        go func(blk_begin_idx int){
            blk_end_idx := blk_begin_idx + partLength

            blk_begin_idx += i_low
            blk_end_idx += i_low

            for i := blk_begin_idx; i < blk_end_idx && i < i_up; i++ {
                for_proc(i)
            }

            ch <- blk_begin_idx
        }(idx)
    }
    for idx := 0; idx < N; idx += partLength {
        <- ch
    }
}
