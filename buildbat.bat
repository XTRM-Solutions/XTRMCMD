go build xtrmcmd
del /F %GOBIN%\xtrmcmd.exe
xcopy xtrmcmd.exe %GOBIN%
