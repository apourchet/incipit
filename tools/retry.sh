#! /bin/bash

CMD=$1
TIME=2
[ -n "$2" ] && TIME=$2
EC=`$CMD`

while [ "$?" -ne "0" ]; do
    sleep $TIME
    echo "Retrying $CMD"
    EC=`$CMD`
done
