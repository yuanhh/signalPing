host="127.0.0.1"
port=7788
highVideoW=1280
hightVideoH=720
#lowVideoW=320
#lowVideoH=180
#lowVideoH=240

while [ $# -gt 0 ]
do
    case "$1" in
    -h)  host="$2"; shift;;
    -p)  port="$2"; shift;;
    --) shift; break;;
    -*)
        echo >&2 \
        "usage: $0 [-d deviceIdx] [-h host] [-p port]"
        exit 1;;
    *)  break;; # terminate while loop
    esac
    shift
done

# capsLow="application/x-rtp, media=(string)video, clock-rate=(int)90000, encoding-name=(string)RAW, sampling=(string)YCbCr-4:2:0, depth=(string)8, width=(string)${lowVideoW}, height=(string)${lowVideoH}, colorimetry=(string)SMPTE240M, payload=(int)96"
capsHigh="application/x-rtp, media=(string)video, clock-rate=(int)90000, encoding-name=(string)RAW, sampling=(string)YCbCr-4:2:0, depth=(string)8, width=(string)${highVideoW}, height=(string)${hightVideoH}, colorimetry=(string)SMPTE240M, payload=(int)96"
# capsAudio="application/x-rtp, media=(string)audio, clock-rate=(int)44100, encoding-name=(string)L16, channels=(int)2, payload=(int)96"
# videoPortLow=$port
videoPortHigh=$((10#$port + 1))
# audioPort=$((10#$port + 2))

gst-launch-1.0 udpsrc port=$videoPortHigh caps="${capsHigh}" ! rtpvrawdepay ! videoconvert ! xvimagesink sync=false
