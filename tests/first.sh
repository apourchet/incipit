#! /bin/bash

echo "Testing hellogo"
curl https://dummy.machine:30002/hellogo -k

echo -e "\nTesting hellonode"
curl https://dummy.machine:30002/hellonode -k
