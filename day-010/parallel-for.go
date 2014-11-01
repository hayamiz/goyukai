//usr/bin/env go run $0 $@ ; exit

package main

import (
    "fmt"
    "time"
    "runtime"
    "math/rand"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    input := make([]int,  10000000)
    output := make([]int, len(input))
    for i, _ := range input {
        input[i] = int(rand.Int31())
    }

    seq_time := 0.0
    pat1_time := 0.0
    pat2_time := 0.0

    for_proc := func(i int) {
        output[i] = input[i] * input[i]
    }

    for i := 0; i < 5; i++ {
        for i, _ := range output {
            output[i] = 0
        }

        t0 := time.Now()
        seq_for(0, len(input) - 1, for_proc)
        dt := time.Since(t0)
        seq_time += dt.Seconds()

        // sanity check
        if i == 0 {
            for i := 0; i < len(input); i++ {
                if output[i] != input[i] * input[i] {
                    panic(i)
                }
            }
        }
        for i, _ := range output {
            output[i] = 0
        }

        t0 = time.Now()
        parallel_for_pat1(0, len(input) - 1, for_proc)
        dt = time.Since(t0)
        pat1_time += dt.Seconds()

        // sanity check
        if i == 0 {
            for i := 0; i < len(input); i++ {
                if output[i] != input[i] * input[i] {
                    panic(i)
                }
            }
        }
        for i, _ := range output {
            output[i] = 0
        }


        t0 = time.Now()
        parallel_for_pat2(0, len(input) - 1, for_proc)
        dt = time.Since(t0)
        pat2_time += dt.Seconds()


        // sanity check
        if i == 0 {
            for i := 0; i < len(input); i++ {
                if output[i] != input[i] * input[i] {
                    panic(i)
                }
            }
        }
    }

    seq_time /= 5.0
    pat1_time /= 5.0
    pat2_time /= 5.0

    fmt.Println("seq time:", seq_time)
    fmt.Println("pat1 time:", pat1_time)
    fmt.Println("pat2 time:", pat2_time)
}

func seq_for(i_low int, i_up int, for_proc func(int)) {
    for i := i_low; i <= i_up; i++ {
        for_proc(i)
    }
}

func parallel_for_pat1(i_low int, i_up int, for_proc func(int)) {
    ch := make(chan int, i_up - i_low + 1)
    for i := i_low; i <= i_up; i++ {
        go func(i int) {
            for_proc(i)
            ch <- i
        }(i)
    }
    for i := i_low; i <= i_up; i++ {
        <- ch
    }
}


func parallel_for_pat2(i_low int, i_up int, for_proc func(int)) {
    ch := make(chan int, i_up - i_low + 1)

    N := i_up - i_low + 1
    blockSize := N / runtime.NumCPU()
    blockNum := runtime.NumCPU() + 1

    ch = make(chan int, blockNum)
    for bi := 0; bi < blockNum; bi++ {
        go func(bi int){
            b_i_low := bi * blockSize
            b_i_up := (bi+1) * blockSize - 1
            if b_i_up >= N - 1 {
                b_i_up = N - 1
            }

            for i := i_low + b_i_low; i <= i_low + b_i_up; i++ {
                for_proc(i)
            }

            ch <- bi
        }(bi)
    }
    for bi := 0; bi < blockNum; bi++ {
        <- ch
    }
}
