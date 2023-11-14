#!/usr/bin/bash

STARTID=${1:-0}
ENDID=${2:-5}
DM_APIGW_PORT=${3:-8000}

function runcmd() {
    local name="${1:-edge-0}"
    #local mydel="${2:-4}"
    #bash -c "sleep $mydel ; echo done-$name"
    bash -c "DM_APIGW_PORT=$DM_APIGW_PORT EDGE_NAME=$name ./bin/eps-edge-agent"
}
function cleanup() {
    #echo `jobs -p`
    #kill ${pids}
    PJOB=$(jobs -p)
    for pid in $PJOB ; do
        kill -9 "$pid" > /dev/null 2> /dev/null || :
    done
}

pids=""
runcmd edge-$STARTID >& /tmp/edge-$STARTID.log&
STARTID=$((STARTID+1))
pids="$pids $!"
sleep 5
for ((i=$STARTID; i < $ENDID ; i++)) ; do 
    runcmd edge-$i >& /tmp/edge-$i.log&
    pids="$pids $!"
done 
echo "PIDs = ${pids}"

trap cleanup SIGINT SIGTERM 
wait $pids
