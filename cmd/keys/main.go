package main

import "github.com/kvvPro/gophkeeper/tools/certs"

func main() {
	certs.MakeRSACert(&certs.Settings{PathToCert: "/workspaces/gophkeeper/cmd/keys/key.pub", PathToPrivateKey: "/workspaces/gophkeeper/cmd/keys/key"})
}
