#! /bin/bash

packages=`go list ../../lib/...`
i=0
for package in $packages; do 
    base=`basename $package`$i
    go test -i -c -o bin/$base $package -args "-v" && echo $package' => ./bin/'$base &
    i=$((i+1))
done

wait
echo "DONE"
