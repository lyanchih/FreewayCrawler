#!/bin/bash

if [ -f '201410261415.csv' ];
then
    mv 201410261415.csv 1414304329.csv
fi

from=1414304329.csv
to=$(date -r $(echo -en $from | awk '{print int($0)}') +"%Y%m%d%H%M.csv" | awk '{printf "%02d.csv\n", int(int($0)/100)*100+int(int($0)%100/5)*5}')

echo -en $from | sed '
h
s/'$from'/'$to'/
x
G
y/\n/ /
s/\(.*\)/mv \1/' | sh

ls *.csv

#sed -i '.bat' $'/,$/s/,$/\\\r/' data/*.csv
