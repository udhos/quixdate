#! /bin/bash

set DEVEL=c:\tmp\devel
set QUIXDATE=%DEVEL%\quixdate
set GOPATH=%QUIXDATE%

[ -d $DEVEL ] || mkdir -p $DEVEL
cd $DEVEL

git clone https://github.com/udhos/quixdate

go_get () {
	local i=$1
	echo go get $i
	go get $i
}

#go_get golang.org/x/net/ipv4