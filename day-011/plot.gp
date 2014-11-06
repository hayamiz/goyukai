set term postscript enhanced color eps

set size 0.6, 0.6

set logscale x
set logscale y

set ylabel "Execution time [sec]"
set xlabel "Input array size"

set key right bottom

set output "| epstopdf -f -o prefixsum.pdf"
plot 'time.dat' usi 1:2 wi lp ti 'sequential prefixsum', \
     'time.dat' usi 1:3 wi lp ti 'paralell prefixsum'


unset logscale y
set yrange [0:*]
set ylabel "Speedup ratio"
set output "| epstopdf -f -o speedup.pdf"
plot 'time.dat' usi 1:($2/$3) wi lp noti
