#!/bin/bash

set -e

while getopts ":d:s:o:" flag
do
    case $flag in
        d)
            FOLDER=$OPTARG
            ;;
        s)
            SORT_COL=$OPTARG
            ;;
        o)
            OUTFILE=$OPTARG
            ;;
        ?)
            echo "Invalid option: -$OPTARG"
            exit -1
            ;;
    esac
done

if [[ -n "$SORT_COL" ]]; then
    cat ${FOLDER:-.}/*.csv | grep -v freeway_id | sort -t, -nk ${SORT_COL}
else
    cat ${FOLDER:-.}/*.csv | grep -v freeway_id
fi
