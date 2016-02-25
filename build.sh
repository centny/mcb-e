#!/bin/bash
set -e
o_dir=build
if [ "$2" != "" ];then
	o_dir=$2
fi
p_dir=$o_dir/mcb-e
rm -rf $o_dir
mkdir -p $o_dir
mkdir -p $p_dir

go get github.com/Centny/ffcm/ffcm
go build -o $p_dir/ffcm github.com/Centny/ffcm/ffcm
#
cp -f run_*.sh $p_dir
cp -f *.properties $p_dir
cp -rf test $p_dir/
cd $o_dir
zip -r mcb-e.zip mcb-e
cd ../
