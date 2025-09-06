success=0
failure=0
error=0
total=0
ftime=15s
failnames=()
passnames=()
errnames=()
onlyone=0
onlyonename=none
fuzzout=none
fuzzexit=999
fuzzid=0
timeblock=0
firstblock=0
mintime=0
while getopts "n:t:d:" opt; do
    case "$opt" in
    n)  onlyone=1
        onlyonename=$OPTARG
        ;;
    t)  ftime=$OPTARG
        ;;
    d)  cd $OPTARG
        ;;
    esac
done
runfuzz() {
    fuzzout=$(go test -fuzz=$1 -fuzztime=$ftime) &
    fuzzid=$!
    printf "┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┨     :    :    :    :    :    :    :    :    :    :    :    :┃"
    printf "\033[1G"
    printf "$1"
    printf "\033[1G\033[40C"
    timeblock=0
    firstblock=0
    mintime=0
    while ps -p "$fuzzid" >/dev/null; do
        if [ "$timeblock" = "0" ]; then
            sleep 1
            if [ "$firstblock" = "0" ]; then
                firstblock=1
                printf "▍"
            else
                printf "\033[1D▒▍"
            fi
        elif [ "$timeblock" = "1" ]; then
            sleep 1
            printf "\033[1D▊"
            mintime=1
        elif [ "$timeblock" = "2" ]; then
            sleep 1
            printf "\033[1D█"
        fi
        let "timeblock = timeblock + 1"
        if [ "$timeblock" = "3" ]; then
            timeblock=0
        fi
    done
    wait $fuzzid
    fuzzexit=$?
    if [ "$fuzzexit" = "0" ]; then
        if [ "$mintime" = "0" ]; then
            printf "\033[1G\033[33m┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈\033[0m"
            printf "\033[1G\033[33m$1\033[0m"
            let "error = error + 1"
            errnames+=($1)
        else
            printf "\033[1G\033[32m┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈\033[0m"
            printf "\033[1G\033[32m$1\033[0m"
            let "success = success + 1"
            passnames+=($1)
        fi
    else
        printf "\033[1G\033[31m┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈\033[0m"
        printf "\033[1G\033[31m$1\033[0m"
        let "failure = failure + 1"
        failnames+=($1)
    fi
    let "total = total + 1"
    echo ""
}
# TESTS
echo "~~~~~~BEGINNING FUZZ TESTS~~~~~~   TIME:    15s  30s  45s  60s  75s  90s  105s 120s 135s 150s 165s 180s"
printf "\033[?25l"
if [ "$onlyone" = "1" ]; then
    runfuzz $onlyonename
else
    runfuzz Fuzz_InsertionSort_
    runfuzz Fuzz_SortedInsert_
    runfuzz Fuzz_SortedSearch_
    cd implementation_test
    runfuzz Fuzz_SliceAdapter_
    runfuzz Fuzz_SliceAdapterIndirect_
    runfuzz Fuzz_FileAdapter_
fi
echo "~~~~~~FUZZ TESTS COMPLETE~~~~~~    TIME:    15s  30s  45s  60s  75s  90s  105s 120s 135s 150s 165s 180s"
# RESULTS
if [ "$success" -gt "0" ]; then
    printf "\033[32m"
    echo "SUCCESS: $success / $total"
    for i in ${!passnames[@]}; do
        echo "  ✓ ${passnames[$i]}"
    done
    printf "\033[0m"
fi
if [ "$error" -gt "0" ]; then
    printf "\033[33m"
    echo "ERRORS:  $error / $total"
    for i in ${!errnames[@]}; do
        echo "  ? ${errnames[$i]} : test may not be located in the current working directory, use -d to specify directory"
    done
    printf "\033[0m"
fi
if [ "$failure" -gt "0" ]; then
    printf "\033[31m"
    echo "FAILURE: $failure / $total"
    for i in ${!failnames[@]}; do
        echo "  X ${failnames[$i]}"
    done
    printf "\033[0m"
fi
if [ "$success" = "$total" ]; then
    printf "\033[0;32mALL TESTS PASS!\033[0m"
    echo ""
fi
printf "\033[?25h"