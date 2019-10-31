
go build -ldflags "-s -w"
cd client
go build -ldflags "-s -w"
cd rev
go build -ldflags "-s -w"
cd ../..

echo "ok!"