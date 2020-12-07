#!/usr/bin/env bash
CSS_FILES=()
while IFS=  read -r -d $'\0'; do
    CSS_FILES+=("$REPLY")
done < <(find $1 -name "*.css" -print0)

for file in "${CSS_FILES[@]}"; do
    MINIFIED=$(cat $file \
        | sed -r ':a; s%(.*)/\*.*\*/%\1%; ta; /\/\*/ !b; N; ba' \
        | tr -d '\t' | tr -d ' ' | tr -d '\n' | tr -s ' ' ' ')
    echo $MINIFIED > $file
done