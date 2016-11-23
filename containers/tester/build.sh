#! /bin/bash

# packages=`go list ../../lib/...`
# i=0
# for package in $packages; do 
#     base=`basename $package`$i
#     echo $package' ... ./bin/'$base && go test -i -c -o bin/$base $package && echo $package' => ./bin/'$base
#     i=$((i+1))
# done
# 
# wait
echo "DONE"
