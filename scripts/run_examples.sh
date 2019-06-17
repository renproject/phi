#!/bin/bash

EXAMPLES_PATH=./task/examples

for dir in $EXAMPLES_PATH/*/
do
		go run ${dir}*
done
