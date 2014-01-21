#!/bin/bash

for d in */
do
	(cd printandexit && ./test.sh) || exit 1
done
