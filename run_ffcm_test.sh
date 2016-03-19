#!/bin/bash
chk_v(){
	tc=`ls -l $1* 2>/dev/null| wc -l`
	if [ $tc -ge $2 ];then
		echo test $3 success
	else
		echo waiting $3 result...
		sleep 3
		chk_v $1 $2 $3
	fi
}
#
tip=$1
if [ "$tip" = "" ];then
	tip=127.0.0.1
fi
echo running test to server $tip
#
echo Test Video...
rm -rf sdata_o/test/xx*
./ffcm -g http://127.0.0.1:2325/addTask?args=test/xx.mp4
chk_v sdata_o/test/xx/ 2 video
echo
echo
#
echo Test docx...
rm -rf sdata_o/test/docx*
./ffcm -g http://127.0.0.1:2325/addTask?args=test/xx.docx,test/docx
chk_v sdata_o/test/docx/d 7 docx
echo
echo
#
echo Test pdfx...
rm -rf sdata_o/test/pdfx*
./ffcm -g http://127.0.0.1:2325/addTask?args=test/xx.pdf,test/pdfx
chk_v sdata_o/test/pdfx/d 6 pdfx
echo
echo
#
echo Test xlsx...
rm -rf sdata_o/test/xlsx*
./ffcm -g http://127.0.0.1:2325/addTask?args=test/xx.xlsx,test/xlsx
chk_v sdata_o/test/xlsx/d 22 xlsx
echo
echo
#
echo Test pptx...
rm -rf sdata_o/test/pptx*
./ffcm -g http://127.0.0.1:2325/addTask?args=test/xx.pptx,test/pptx
chk_v sdata_o/test/pptx/d 1 pptx
echo
echo
#
echo Test png...
rm -rf sdata_o/test/png*
./ffcm -g http://127.0.0.1:2325/addTask?args=test/xx.jpg,test/png
chk_v sdata_o/test/png 1 png
echo
echo
#
echo all test done...