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
OS=`uname -a`

function unixToTime() {
    if [[ $OS == Darwin* ]];
    then
        date -u -r $1 +"%Y%m%d%H%M" | awk '{printf "%d%02d", int(int($0)/100), int(int($0)%100/5)*5}'
    else
        date -u --date="@$1" +"%Y%m%d%H%M" | awk '{printf "%d%02d", int(int($0)/100), int(int($0)%100/5)*5}'
    fi
}

function defaultFromTime() {
    if [[ $OS == Darwin* ]];
    then
        date -v-1H +"%s"
    else
        date -d '1 hour ago' +"%s"
    fi
}

function parseParameter() {
    local request=$1
    
    for para in $(echo -en ${request#/*\?} | sed 's/&/ /g');
    do
        local key=${para%%=*} value=${para##*=}
        case $key in
            from)
                if [ ${#value} == 10 ];
                then
                    from=$value
                fi
                ;;
            to)
                if [ ${#value} == 10 ];
                then
                    to=$value
                fi
                ;;
        esac
    done
    if [ ${#from} == 0 ];
    then
        if [ ${#to} == 0 ];
        then
            from=$(defaultFromTime)
        else
            from=$((to-3600))
        fi
    fi
    
    if [ ${#to} == 0 ];
    then
        to=$((from+3600))
    fi
    
    from=$(unixToTime $from)
    to=$(unixToTime $to)
}

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
    unset from
    unset to
    
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
    echo $(date): $REQUEST
    
    if [ "$REQUEST" == "/highway.json" -o "$REQUEST" == "/highway.min.json" ];
    then
        file=${REQUEST##*/}
        echo -en "HTTP/1.1 200 OK\r\nAccess-Control-Allow-Origin: *\r\nContent-Length: $(ls -la $file | awk '{print $5}')\r\nContent-Type: application/json\r\n\r\n" | cat - $file >$out
        return
    elif [ "$REQUEST" == "/highway.min.json.gz" ];
    then
        file=${REQUEST##*/}
        echo -en "HTTP/1.1 200 OK\r\nAccess-Control-Allow-Origin: *\r\nContent-Length: $(ls -la $file | awk '{print $5}')\r\nContent-Encoding: gzip\r\nContent-Type: application/json\r\n\r\n" | cat - $file >$out
        return
    elif [ "${REQUEST#/api}" == "$REQUEST" ];
    then
        echo -en "HTTP/1.1 404 OK\r\n\r\n" >$out
        return
    fi
    
    parseParameter $REQUEST
    
    fs=$(ls -a $FOLDER/*.csv | sed -E 's/^.*\/([[:digit:]]+)\.csv/\1/' | awk -v from=$from -v to=$to -v folder=$FOLDER '$1 >= from && $1 < to { printf "%s/%s.csv\n", folder, $1}')
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
