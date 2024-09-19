#!/bin/bash

set -e

if [ $# != 1 ] ; then
    echo "USAGE: $0 commitMessage"
    echo "
            e.g.: $0 \"commit message\"
        "
    exit 1;
fi

find . -name .DS_Store | xargs rm -rf

rm -rf scripts/java/{logs,temp}

git add . -A
git commit -m "$1"
git push
