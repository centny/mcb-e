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
rm -rf out/test/xx_*
./ffcm -g http://127.0.0.1:2325/v/addTask?src=test/xx.mp4
chk_v out/test/xx_ 2 video
echo
echo
#
echo Test docx...
rm -rf out/test/docx_*
./ffcm -g http://127.0.0.1:2325/n/addTask?args=test/xx.docx,test/docx
chk_v out/test/docx_ 7 docx
echo
echo
#
echo Test pdfx...
rm -rf out/test/pdfx_*
./ffcm -g http://127.0.0.1:2325/n/addTask?args=test/xx.pdf,test/pdfx
chk_v out/test/pdfx_ 6 pdfx
echo
echo
#
echo Test xlsx...
rm -rf out/test/xlsx_*
./ffcm -g http://127.0.0.1:2325/n/addTask?args=test/xx.xlsx,test/xlsx
chk_v out/test/xlsx_ 22 xlsx
echo
echo
#
echo Test pptx...
rm -rf out/test/pptx_*
./ffcm -g http://127.0.0.1:2325/n/addTask?args=test/xx.pptx,test/pptx
chk_v out/test/pptx_ 1 pptx
echo
echo
#
echo all test done...