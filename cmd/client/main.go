package main

import "github.com/kvvPro/gophkeeper/cmd/client/cli"

// go build -ldflags "-X 'github.com/kvvPro/gophkeeper/cmd/client/cli.BuildVersion=v1.0.1' -X 'github.com/kvvPro/gophkeeper/cmd/client/cli.BuildDate=$(date +'%Y/%m/%d %H:%M:%S')'"

func main() {
	cli.Execute()
}
