@echo off
setlocal

rem Define array elements
set "PROTO_NAMES=user"

rem Loop through each element in the array
for %%i in (%PROTO_NAMES%) do (
    protoc --go_out=plugins=grpc:./%%i --go_opt=module=github.com/openimsdk/openmeeting-server/pkg/protocol/%%i %%i/%%i.proto
    if ERRORLEVEL 1 (
        echo error processing %%i.proto
        exit /b %ERRORLEVEL%
    )
)

rem Replace "omitempty" in *.pb.go files with UTF-8 encoding
for /r %%f in (*.pb.go) do (
    powershell -Command "(Get-Content -Path '%%f' -Encoding UTF8) -replace ',omitempty\"`"', '\"`"' | Set-Content -Path '%%f' -Encoding UTF8"
)


endlocal
