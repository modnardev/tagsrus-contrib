#!/bin/bash

HOST="http://127.0.0.1:8080"
AFTER=""

get () {
    out=$(curl -s -A 'User-Agent: tagsrus-sampler/1.0.0' \
    "https://www.reddit.com/r/all/top.json?t=all&after=$AFTER" \
    | jq -r '.data.children[]|.data.name')

    for i in $out; do
        curl -s "$HOST/api/import?link=https://reddit.com/${i:3}" > /dev/null
        echo ${i:3}
        AFTER=$i
    done
}

for i in $(seq 1 1 20); do
    get
done
