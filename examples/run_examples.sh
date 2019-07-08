#!/bin/bash

EXAMPLES_PATH=./examples

for dir in $EXAMPLES_PATH/*/
do
		go run ${dir}*
		if [ $? != 0 ]; then
				exit 1
		fi
done
