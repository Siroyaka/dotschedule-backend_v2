YDAY=`date "+%Y/%-m/%-d" -d yesterday`

DURL="https://dotschedule.siroyaka.net/schedule/$YDAY"

DISCORDURL=""

curl -X POST -H 'Content-Type: application/json' -d "{\"content\" : \"$YDAYのどっとライブYoutube配信情報\n$DURL\"}" $DISCORDURL
