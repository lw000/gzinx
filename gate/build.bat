set prjPath=%cd%
echo %prjPath%
cd ../../../../
set GOPATH=%cd%
set GOARCH=amd64
set GOOS=windows
cd %prjPath%
rem go mod vendor
go build -v -ldflags="-s -w"
