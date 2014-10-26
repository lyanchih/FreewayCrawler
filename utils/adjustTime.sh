#!/bin/bash

for file in ./data/*.csv;
do
    sed -i '.bat' -E '
2btime
3,$ {s/^[[:digit:]]+,/,/; x; G; s/\n//; btime
}

:time
h; s/^([[:digit:]]+),.*$/\1/; x; b
' $file
done
