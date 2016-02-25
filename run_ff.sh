#!/bin/bash
set -e
export PATH=`pwd`:`dirname ${0}`:$PATH
mkdir -p `dirname $7`
mkdir -p `dirname $8`
tsize_=`ffcm -d $3 $4 $5 $6`
ffmpeg -progress $1 -i $2 -s $tsize_ -y $7
cp -f $7 $8
rm -f $7
echo
echo
echo "----------------result----------------"
echo "[text]"
echo $9
echo "[/text]"
echo