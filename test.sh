#!/usr/bin/env bash

function unit_test() {
    local input=$1
    local expected=$2
    echo "$input" | go run main.go > a.s
    gcc a.s
    ./a.out
    local actual=$?
    if [[ $expected -eq $actual ]];then
        echo "ok"
    else
        echo "not ok : $expected != $actual"
        exit 1
    fi
}

unit_test '42' '42'
unit_test '+7' '7'
unit_test ' 7' '7'
unit_test '7;' '7'
unit_test ' 7 ;' '7'
unit_test '-1' '255'
unit_test '30+12' '42'
unit_test '6 * 7' '42'
unit_test '42 / 2' '21'

