#!/bin/bash

# HOST="http://192.168.1.100:8888" # sfw
HOST="http://192.168.1.100:8889" #nsfw
# HOST="http://localhost:8090" #mixed
MAX=500

PAGES=( \
# "https://www.reddit.com/r/all" \
)

# hour, day, week, month, year, all
T="month"

get () {
    LINK="$PAGE/top.json?t=$T&after=$AFTER"
    echo "-- $LINK"
    out=$(curl -s -A 'User-Agent: tagsrus-sampler/1.0.0' $LINK \
    | jq -r '.data.children[]|.data.name+";"+.data.permalink')

    for i in $out; do
        j=(${i//;/ })
        echo ${j[0]:3} ${j[1]}
        curl -s "$HOST/api/import?link=https://www.reddit.com${j[1]}" > /dev/null
        AFTER=${j[0]}
    done
}

for PAGE in ${PAGES[@]}; do
    AFTER=""
    for i in $(seq 1 25 $(($MAX-1))); do
        get $PAGE
    done
done
