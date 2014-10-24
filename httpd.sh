#!/bin/bash

out=/tmp/pipe

rm -f $out
mkfifo $out
trap "rm -f /tmp/pipe" EXIT

while getopts ":d:" flag
do
    case $flag in
        d)
            FOLDER=$OPTARG
            ;;
    esac
done

ls ${FOLDER:=./data}/*.csv >/dev/null 2>&1
if [ "$?" != "0" ];
then
    echo "$FOLDER don't have any csv file"
    exit 1
fi

HEADER=$(awk 'NR = 1 {print; exit; }' $(ls -a ./data/*.csv | head -n1))"\n"

function parse_header {
    x=0
    while read line[$x] && [ ${#line[$x]} -gt 1 ];
    do
        x=$(($x+1))
    done
    
    unset line[$((${#line[@]}-1))]
    METHOD=$(echo ${line[0]} |cut -d" " -f1)
    REQUEST=$(echo ${line[0]} |cut -d" " -f2)
    HTTP_VERSION=$(echo ${line[0]} |cut -d" " -f3)
    local $(echo ${REQUEST#/*\?} | sed 's/&/ /g')
    
    from="$FOLDER/${from:=$(date +"%s")}.csv"
    to="$FOLDER/${to:=$(date +"%s")}.csv"
    fs=$(ls -a $FOLDER/*.csv | awk -v from=$from -v to=$to '$1 >= from && $1 <= to { print $1} ' | sort)
    if [ -z "$fs" ];
    then
        body=$HEADER
    else
        body=$HEADER$(awk 'FNR>1{print;}' $fs)
    fi
    echo -en "HTTP/1.1 200 OK\r\nAccess-Control-Allow-Origin: http://null.jsbin.com\r\nContent-Length: $(echo -en $body | wc -c)\r\nContent-Type: text/csv\r\n\r\n$body">$out
}

while true
do
    cat $out | nc -l 8080 | parse_header
done
