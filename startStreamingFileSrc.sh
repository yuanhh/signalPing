filePath=./Simpsons.mp4
port=7788
highVideoW=1280
hightVideoH=720
#lowVideoW=320
#lowVideoH=240
#debug=

#while [ $# -gt 0 ]
#do
#    case "$1" in
#    -f)  filePath="$2"; shift;;
#    -h)  host="$2"; shift;;
#    -p)  port="$2"; shift;;
#    -d)  debug=1;;
#    --) shift; break;;
#    -*)
#        echo >&2 \
#        "usage: $0 [-f filepath] [-h host] [-p port] [-d]"
#        exit 1;;
#    *)  break;; # terminate while loop
#    esac
#    shift
#done

echo "launching signalServer"

if [[ -x "signalServer" ]];
then
    host=$(sudo ./signalServer)
    echo $host
else
    echo "signalServer not found"
fi

#capsLow="video/x-raw,width=${lowVideoW},height=${lowVideoH}"
capsHigh="video/x-raw,width=${highVideoW},height=${hightVideoH}"
#videoPortLow=$port
videoPortHigh=$port
audioPort=$((10#$port + 1))
#echo "listening on ${videoPortLow} for low video resolution"
echo "listening on ${videoPortHigh} for high video resolution"
echo "listening on $audioPort for audio"

while true
do
if [ -z $debug ]
then
gst-launch-1.0 -v filesrc location="${filePath}" ! decodebin name=dec ! tee name=t \
t. ! queue ! videoconvert ! videoscale ! $capsHigh ! rtpvrawpay ! udpsink host=$host  port=$videoPortHigh \
dec. ! audioresample ! audioconvert ! audio/x-raw,channels=2,depth=16,width=16,rate=44100 ! rtpL16pay ! udpsink host=$host port=$audioPort 2>&1 > /dev/null
else
gst-launch-1.0 -v filesrc location="${filePath}" ! decodebin name=dec ! tee name=t \
t. ! queue ! videoconvert ! videoscale ! $capsHigh ! rtpvrawpay ! udpsink host=$host  port=$videoPortHigh \
dec. ! audioresample ! audioconvert ! audio/x-raw,channels=2,depth=16,width=16,rate=44100 ! rtpL16pay ! udpsink host=$host port=$audioPort
fi
done
