#the power shell
if($args.Length -lt 9){
	echo "bad arguments"
	exit 1
}
Write-Host 'run_ff arguments list:'
echo '  '$args[0]
echo '  '$args[1]
echo '  '$args[2]
echo '  '$args[3]
echo '  '$args[4]
echo '  '$args[5]
echo '  '$args[6]
echo '  '$args[7]
echo '  '$args[8]
echo
$dst_d=Split-Path -Parent $args[6]
$tmp_d=Split-Path -Parent $args[7]
mkdir -p $dst_d
mkdir -p $tmp_d
tsize_=ffcm -d $args[2] $args[3] $args[4] $args[5]
ffmpeg -progress $args[0] -i $args[1] -s $tsize_ -y $args[6]
echo 'do copy '+$args[6]+' to '+$args[7]
copy -f $args[6] $args[7]
echo 'do remove tmp file '+$args[6]
rm -f $args[6]
echo 'all done...'
echo
echo
echo "----------------result----------------"
echo "[text]"
echo $args[8]
echo "[/text]"
echo
echo
sleep 1
echo
echo