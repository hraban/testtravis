#!/bin/bash

for d in */
do
	echo testing $d
	(cd printandexit && ./test.sh) || exit 1
done
