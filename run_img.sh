#!/bin/bash
set -e
export PATH=`pwd`:`dirname ${0}`:$PATH
mkdir -p `dirname $3`
mkdir -p `dirname $4`
args=""
if [ "$2" != "" ];then
	args="-resize "$2
fi
convert $1 $args $3
cp -f $3 $4
rm -f $3
echo
echo
echo "----------------result----------------"
echo "[text]"
echo $5
echo "[/text]"
echo