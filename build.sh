#!/bin/bash
set -e
if [ "$1" == "" ];then
	echo "Usage: ./build.sh <linux|win> <cswf.doc location> <cswf.srv location>"
	exit 1
fi
o_dir=build
pkgn=mcb-e.$1
p_dir=$o_dir/$pkgn
rm -rf $o_dir
mkdir -p $o_dir
mkdir -p $p_dir

sys_n=`uname`
out_n=ffcm
if [ ${sys_n:0:7} = "MSYS_NT" ];then
	out_n=ffcm.exe
fi
go get github.com/Centny/ffcm/ffcm
go build -o $p_dir/$out_n github.com/Centny/ffcm/ffcm
#
cp -f run_*.sh $p_dir
cp -f run_*.bat $p_dir
cp -f do_* $p_dir
cp -f *.properties $p_dir
cp -rf test $p_dir/
if [ "$2" != "" ];then
	cp -rf $2/* $p_dir
fi
cd $o_dir
zip -r $pkgn.zip $pkgn
cd ../
