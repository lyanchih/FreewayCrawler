#!/bin/bash

while getopts ":ut" flag;
do
    case $flag in
        u)
            if [ "$type" != "" ];
            then
                echo "You must specify -u or -t"
                exit 1
            fi
            type="u"
            ;;
        t)
            if [ "$type" != "" ];
            then
                echo "You must specify -u or -t"
                exit 1
            fi
            type="t"
            ;;
        *)
            echo "not allow flag " $OPTARG
            exit 1
            ;;
    esac
done
if [ "$type" == "" ];
then
    echo "You must specify -u or -t"
    exit 1
fi

for file in ./data/*.csv;
do
    from=$file
    if [ "$type" == "t" ];
    then
        to=$(echo $from | sed -E '/\.csv$/s/.*\/([^/]+)\.csv$/date -r \1 +"%Y%m%d%H%M"/' | bash | awk '{printf "%d.csv", int(int($0)/100)*100+int(int($0)%100/5)*5}' | cat - <(echo $from | sed -E '/\.csv$/s/^(.*\/)[^/]+\.csv$/ \1/') | awk '{printf "%s%s", $2, $1}')
    else [ "$type" == "u" ];
        to=$(echo $from | sed -E '/\.csv$/s/.*\/([^/]+)\.csv$/date -jf "%Y%m%d%H%M%S" \100 +"%s"/' | bash | awk '{printf "%s.csv", $0}' | cat - <(echo $from | sed -E '/\.csv$/s/^(.*\/)[^/]+\.csv$/ \1/') | awk '{printf "%s%s", $2, $1}')
    fi
    mv $from $to
done

# if [ -f '201410261415.csv' ];
# then
#     mv 201410261415.csv 1414304329.csv
# fi

# from=1414304329.csv
# to=$(date -r $(echo -en $from | awk '{print int($0)}') +"%Y%m%d%H%M.csv" | awk '{printf "%02d.csv\n", int(int($0)/100)*100+int(int($0)%100/5)*5}')

# echo -en $from | sed '
# h
# s/'$from'/'$to'/
# x
# G
# y/\n/ /
# s/\(.*\)/mv \1/' | sh

# ls *.csv
