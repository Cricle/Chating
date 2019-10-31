
$oldcgo=$env:CGO_ENABLED
$oldgoos=$env:GOOS
$oldgoarch=$env:GOARCH
$env:GOOS="linux"
$env:GOARCH="amd64"


go build -ldflags "-s -w"


$env:GOOS=$oldgoos
$env:GOARCH=$oldgoarch
echo "ok!"