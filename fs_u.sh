#!/bin/bash
for file in $1;
do
	echo upload $file
	gfs -u http://127.0.0.1:2010 $file "" $2
done
