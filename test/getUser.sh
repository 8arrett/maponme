#!/bin/sh

### All curl tests moved into api.py

test () {
    if [ -z $3 ]; then
        echo "Test called without full arguments. Server down?"
        exit 1
    fi
    if [ ! $2 = $3 ]; then
        echo "Failed: " $1
        echo $2
        echo $3
        exit 1
    fi
}

resp=$(curl -s localhost:3000/api/error)
test 'Basic test' $resp '{"Status":"Fail"}'

resp=$(curl -s localhost:3000/api/erroragain)
test 'Basic test' $resp '{"Status":"Fail"}'

echo "Passed"
return 0