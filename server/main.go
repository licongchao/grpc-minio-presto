package main

import (
	"os"

	"github.com/urfave/cli"
)

// var (
// 	port = flag.Int("port", 50000, "Model Service port")
// )

func main() {
	app := cli.NewApp()
	app.Name = "Model Operation"
	app.Usage = "Model Version Storage"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		StartServerCommand(),
	}
	app.Run(os.Args)
}
