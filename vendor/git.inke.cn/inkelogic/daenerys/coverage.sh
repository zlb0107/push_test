#!/bin/bash

set -e

profile="cover.out"
htmlfile="cover.html"
mergecover="merge_cover"
mode="set"

for package in $(go list ./...|grep -v src); do
    coverfile="$(echo $package | tr / -).cover"
    go test -covermode="$mode" -coverprofile="$coverfile" -coverpkg="$package" "$package"
done

# merge all profiles
grep -h -v "^mode:" *.cover | sort > $mergecover

# aggregate duplicated code-block data
echo "mode: $mode" > $profile
current=""
count=0
while read line; do
    block=$(echo $line | cut -d ' ' -f1-2)
    num=$(echo $line | cut -d ' ' -f3)
    if [ "$current" == "" ]; then
        current=$block
        count=$num
    elif [ "$block" == "$current" ]; then
        count=$(($count + $num))
    else
        echo $current $count >> $profile
        current=$block
        count=$num
    fi
done < $mergecover

if [ "$current" != "" ]; then
    echo $current $count >> $profile
fi

# save result
# go tool cover -html=$profile -o $htmlfile
go tool cover -func=$profile -o coverage.txt
tail -n 1 coverage.txt|awk -F ' ' '{print $1,$3}'