#!/usr/bin/env bash

CONTENT=$(go run *.go)
if [[ $? -ne 0 ]]
then
  echo "Error:"
  echo "$CONTENT"
  exit 1
fi

echo "" > README.md
while read p; do
  if [[ $p == "{{ CLIENTS }}" ]]
  then
    echo "$CONTENT" >> README.md
  else
    echo "$p" >> README.md
  fi
done < README.template.md