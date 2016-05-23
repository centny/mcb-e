#!/bin/bash
set -e
if [ $# -lt 9 ];then
	echo "bad arguments"
	exit 1
fi
echo run_ff arguments list:
echo '  '$1
echo '  '$2
echo '  '$3
echo '  '$4
echo '  '$5
echo '  '$6
echo '  '$7
echo '  '$8
echo '  '$9
echo
export PATH=`pwd`:`dirname ${0}`:$PATH
mkdir -p `dirname $7`
mkdir -p `dirname $8`
tsize_=`ffcm -d $3 $4 $5 $6`
ffmpeg -progress $1 -i $2 -s $tsize_ -y $7
echo 'do copy '$7' to '$8
cp -f $7 $8
echo 'do remove tmp file '$7
rm -f $7
echo 'all done...'
echo
echo
echo "----------------result----------------"
echo "[text]"
echo $9
echo "[/text]"
echo
echo
sleep 1
echo
echo
