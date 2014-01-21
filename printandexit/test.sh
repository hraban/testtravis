#!/bin/bash

function finish {
	rm -f ./status
}
trap finish EXIT

go build . || exit 1
./printandexit test bar 2> ./status | nonexistingcmd
status=$(< status)
echo  output was $status
