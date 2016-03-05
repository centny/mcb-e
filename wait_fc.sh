#!/bin/bash
for ((i=0; i>-1; ++1))
do
	tc=`ls -l $1 2>/dev/null| wc -l`
	if [ $tc -ge $2 ];then
		echo test $3 success
		break
	else
		echo waiting $3 $2 result...
		sleep 3
	fi
done