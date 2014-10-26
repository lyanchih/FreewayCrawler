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

function handleEncoding() {
    local line=$@
    local formats=$(echo ${line##Accept-Encoding: } | tr ',' ' ')
    for format in $formats
    do
        case ${format%%;*} in
            gzip)
                compress='gzip'
                ;;
        esac
    done
}

function parse_header {
    x=0
    compress=""
    while read line[$x] && [ ${#line[$x]} -gt 1 ];
    do
        if [ -z "${line[$x]##Accept-Encoding: *}" ]; then handleEncoding ${line[$x]}; fi
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
        body=$(awk -v HEADER=$HEADER 'BEGIN {print HEADER} FNR>1 {print; }' $fs)
    fi
    body_len=$(echo -en $body | wc -c)
    if [ -z "$compress" ];
    then
        echo -en "HTTP/1.1 200 OK\r\nAccess-Control-Allow-Origin: *\r\nContent-Length: $body_len\r\nContent-Type: text/csv\r\n\r\n" | cat - <(echo -en $body | tr ' ' "\n")>$out
    else
        body_len=$(echo -en $body | tr ' ' "\n" | gzip -1c | wc -c)
        echo -en "HTTP/1.1 200 OK\r\nAccess-Control-Allow-Origin: *\r\nContent-Length: $body_len\r\nContent-Encoding: gzip\r\nContent-Type: text/csv\r\n\r\n" |
        cat - <(echo -en $body | tr ' ' "\n" | gzip -1c) >$out
    fi
}

while true
do
    cat $out | nc -l 8080 | parse_header
done
