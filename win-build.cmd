@rem win-goinstall

set DEVEL=c:\tmp\devel
set QUIXDATE=%DEVEL%\quixdate
set GOPATH=%QUIXDATE%

go tool fix %QUIXDATE%\src
go tool vet %QUIXDATE%\src
gofmt -s -w %QUIXDATE%\src

@rem build server
go install webd

@rem eof
