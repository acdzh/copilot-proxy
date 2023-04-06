mkdir -p ./output

cd ./client
# macos
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build . && mv ./copilot-proxy ../output/copilot-proxy__darwin64
# # linux
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build . && mv ./copilot-proxy ../output/copilot-proxy__linux64
# # windows
# CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build . && mv ./copilot-proxy.exe ../output/copilot-proxy__win64.exe
cd ..

cd ./server
# macos
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build proxy.go && mv ./proxy ../output/copilot-proxy-server__darwin64
# linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build proxy.go && mv ./proxy ../output/copilot-proxy-server__linux64
# windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build proxy.go && mv ./proxy.exe ../output/copilot-proxy-server__win64.exe