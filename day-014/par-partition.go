//usr/bin/env go run $0 $@ ; exit

package main

import (
    "fmt"
    "runtime"
    "math/rand"
    "time"
    "math"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    TRY := 10

    fmt.Printf("#N	seq_time	par_time\n")

    for N := 1024; N < 1e9; N *= 2 {
        seq_time := 0.0
        par_time := 0.0
        input1 := make([]int,  N)
        input2 := make([]int,  N)
        output := make([]int, len(input1))
        par_output := make([]int, len(input1))
        tmp := make([]int, len(input1))

        for try := 0; try < TRY; try ++ {
            for i, _ := range input {
                input1[i] = int(rand.Int31())
                input2[i] = input1[i]
            }

            var k1, k2 int
            seq_time += exec_time(func() {
                k1 := partition(input1, 0, len(input1))
            })

            par_time += exec_time(func() {
                k2 := par_partition(input2, 0, len(input2))
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

        fmt.Printf("%d	%f	%f\n", N, seq_time, par_time)
    }
}

func exec_time(body_proc func()) float64 {
    t0 := time.Now()
    body_proc()
    dt := time.Since(t0)
    return dt.Seconds()
}

func prefixsum(input []int, output []int) {
    for i, x := range input {
        if i == 0 {
            output[i] = x
        } else {
            output[i] = output[i - 1] + x
        }
    }
}

func partition(input []int, low int, up int) int {
    // sanity check
    if up - low < 0 {
        panic("quick_sort: Invalid argument")
    }

    if up - low <= 1 {
        return
    }

    if up - low == 2 {
        if data[low] > data[up - 1] {
            data[low], data[up - 1] = data[up - 1], data[low]
        }
        return
    }

    pivot := data[low + (up - low) / 2]

    i, j := low, up - 1
    for {
        // move i
        for data[i] < pivot && i < j {
            i++
        }

        // move j
        for data[j] > pivot && i < j {
            j--
        }

        if i < j {
            data[i], data[j] = data[j], data[i]
            i++
            continue
        } else {
            if data[i] != pivot {
                panic(i)
            } else {
                return i
            }
        }
    }
}

func par_prefixsum(input []int, output []int, tmp []int) {
    batchSize := 256

    if len(input) <= 0 {
        prefixsum(input, output)
    } else {
        // y := make([]int, len(input) / 2)
        // z := make([]int, len(input) / 2)
        y := output[0:len(input)/2]
        z := tmp[0:len(input)/2]

        numBatch := int(math.Ceil(float64(len(y)) / float64(batchSize)))

        parallel_for(0, numBatch, func(batch_idx int){
            i_end := (batch_idx + 1) * batchSize
            if i_end > len(y) {
                i_end = len(y)
            }
            for i := batch_idx * batchSize; i < i_end; i++ {
                y[i] = input[2 * i] + input[2 * i + 1]
            }
        })

        par_prefixsum(y, z, tmp[len(input)/2:len(input)/2 + len(input)/2])

        numBatch = int(math.Ceil(float64(len(input)) / float64(batchSize)))
        parallel_for(0, numBatch, func(batch_idx int){
            i_end := (batch_idx + 1) * batchSize
            if i_end > len(input) {
                i_end = len(input)
            }
            for i := batch_idx * batchSize; i < i_end; i++ {
                if i == 0 {
                    output[0] = input[0]
                } else if i % 2 == 1 {
                    output[i] = z[i / 2]
                } else {
                    output[i] = z[(i - 1) / 2] + input[i]
                }
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
