#!/bin/bash

for file in ./data/*.csv;
do
    sed -i '.bat' $'/,$/s/,$/\\\r/' $file
done
