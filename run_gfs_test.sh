#!/bin/bash
#
set -e
tip=$1
if [ "$tip" = "" ];then
	tip=127.0.0.1
fi
echo running test to server $tip
rm -rf out/www/*
#
echo Test Video...
./gfs -u http://127.0.0.1:2325 test/xx.mp4
./wait_fc.sh "out/www/*" 2 video
echo
echo
#
echo Test docx...
./gfs -u http://127.0.0.1:2325 test/xx.docx
./wait_fc.sh "out/www/*" 9 docx
echo
echo
#
echo Test pdfx...
./gfs -u http://127.0.0.1:2325 test/xx.pdf
./wait_fc.sh "out/www/*" 15 pdfx
echo
echo
#
echo Test xlsx...
./gfs -u http://127.0.0.1:2325 test/xx.xlsx
./wait_fc.sh "out/www/*" 37 xlsx
echo
echo
#
echo Test pptx...
./gfs -u http://127.0.0.1:2325 test/xx.pptx
./wait_fc.sh "out/www/*" 38 pptx
echo
echo
#
echo all test done...