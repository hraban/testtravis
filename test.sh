#!/bin/bash

for d in */
do
	echo testing $d
	(cd $d && ./test.sh) || exit 1
done
