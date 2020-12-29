#!/usr/bin/env bash
CSS_FILES=()
while IFS=  read -r -d $'\0'; do
    CSS_FILES+=("$REPLY")
done < <(find $1 -name "*.css" -print0)

for file in "${CSS_FILES[@]}"; do
    MINIFIED=$(cat $file \
        | tr -d '\t' | tr -d '\n' | tr -s ' ' ' ' \
        | sed -r -e 's/([;:{]) ([a-z-])/\1\2/g')
    echo $MINIFIED > $file
done