#!/bin/bash
mkdir -p demo-dir
total=1000000

for ((i=1; i<=total; i++)); do
    : > "demo-dir/$i"
    if (( i % 1000 == 0 )); then
        echo -ne "Progress: $(( i * 100 / total ))% ($i/$total)\r"
    fi
done

echo -e "\nDone!"

