#!/bin/bash
set -e
export PATH=`pwd`:`dirname ${0}`:$PATH
mkdir -p `dirname $3`
mkdir -p `dirname $2`
ffmpeg -i $1 -write_xing 0 $3
cp -f $3 $2
rm -f $3
echo
echo
echo "----------------result----------------"
echo "[json]"
echo '{"count":1,"files":["'$4'"],"src":"'$1'"}'
echo "[/json]"
echo