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

    // test_input := make([]int, 20)
    // test_input2 := make([]int, 20)
    // for i, _ := range test_input {
    //     test_input[i] = int(rand.Intn(100))
    //     test_input2[i] = test_input[i]
    // }
    // par_partition(test_input, 0, len(test_input) - 1, test_input[len(test_input) / 2 - 2])
    // partition(test_input2, 0, len(test_input) - 1, test_input[len(test_input) / 2 - 2])
    // 
    // fmt.Println(test_input2)


    TRY := 10

    fmt.Printf("#N	seq_time	par_time\n")

    for N := 1024; N < 1e9; N *= 2 {
        seq_time := 0.0
        par_time := 0.0
        input := make([]int,  N)
        par_input := make([]int,  N)

        for try := 0; try < TRY; try ++ {
            for i, _ := range input {
                input[i] = int(rand.Int31())
                par_input[i] = input[i]
            }

            k := -1
            par_k := -1
            x, y, z := input[0], input[len(input)/2], input[len(input) - 1]

            // pivot selection
            var pivot int
            if y < x {
                x, y = y, x
            }
            if z < x {
                pivot = x
            } else if y < z {
                pivot = y
            } else {
                pivot = z
            }

            seq_time += exec_time(func() {
                k = partition(input, 0, len(input) - 1, pivot)
            })

            par_time += exec_time(func() {
                par_k = par_partition(input, 0, len(input) - 1, pivot)
            })

            // sanity check
            if k != par_k {
                panic(fmt.Sprintf("Partition index mismatch: k=%d, par_k=%d", k, par_k))
            }
            for j, _ := range input {
                if (j < k && input[j] >= pivot) || (j > k && input[j] <= pivot) {
                    panic("Failed to partition.")
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

// pivot must be in input[q:r]
func partition(input []int, q int, r int, pivot int) int {
    i, j := q, r
    for {
        for input[i] < pivot && i < j {
            i++
        }

        for input[j] > pivot && i < j {
            j--
        }

        if i < j {
            input[i], input[j] = input[j], input[i]
            continue
        } else {
            if input[i] != pivot {
                fmt.Printf("input[i-1]=%d,input[i]=%d,input[i+1]=%d\n", input[i-1], input[i], input[i+1])
                panic(fmt.Sprintf("pivot is not in input: q=%d, r=%d, i=%d, j=%d, pivot=%d",
                    q, r, i, j, pivot))
            }
            return i
        }
    }
}

func par_partition(input []int, q int, r int, pivot int) int {
    // fmt.Println("par_partition input:", input)
    // fmt.Println("par_partition pivot:", pivot)

    n := r - q + 1
    if n == 1 {
        return q
    }

    B, lt, gt := make([]int, n), make([]int, n), make([]int, n)
    parallel_for(0, n, func(i int){
        B[i] = input[q + i]
        if B[i] < pivot {
            lt[i] = 1
        } else {
            lt[i] = 0
        }
        if B[i] > pivot {
            gt[i] = 1
        } else {
            gt[i] = 0
        }
    })

    // fmt.Println(lt)
    // fmt.Println(gt)

    lt_psum, gt_psum, psum_tmp := make([]int, n), make([]int, n), make([]int, n)
    par_prefixsum(lt, lt_psum, psum_tmp)
    par_prefixsum(gt, gt_psum, psum_tmp)

    k := q + lt_psum[n - 1]
    input[k] = pivot

    parallel_for(0, n, func(i int){
        if B[i] < pivot {
            input[q + lt_psum[i] - 1] = B[i]
        } else if B[i] > pivot {
            input[k + gt_psum[i]] = B[i]
        }
    })

    // fmt.Println("par_partition output: ", input)
    // fmt.Println("par_partition index: ", k)

    return k
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
