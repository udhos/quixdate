@rem win-run

set DEVEL=c:\tmp\devel
set QUIXDATE=%DEVEL%\quixdate
set BIN=%QUIXDATE%\bin

start cmd /k %BIN%\webd -config=%QUIXDATE%\config-webd.txt

@rem eof
