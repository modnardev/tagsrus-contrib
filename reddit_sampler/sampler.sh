#!/bin/bash

# HOST="http://192.168.1.100:8888" # sfw
HOST="http://192.168.1.100:8889" #nsfw
# HOST="http://localhost:8090" #mixed
MAX=500

# hour, day, week, month, year, all
T="month"

PAGES=( \
# "https://www.reddit.com/r/all" \
)

get () {
    LINK="$PAGE/top.json?t=$T&after=$AFTExR"
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

if [ -n "$1" ]; then
    T=$1
fi

if [ -n "$2" ]; then
    MAX=$(($2))
fi

for PAGE in ${PAGES[@]}; do
    AFTER=""
    for i in $(seq 1 25 $(($MAX-1))); do
        get $PAGE
    done
done
