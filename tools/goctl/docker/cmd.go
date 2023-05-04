package docker

import "github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"

var (
	varServiceName  string
	varServiceType  string
	varStringBase   string
	varIntPort      int
	varStringHome   string
	varStringRemote string
	varStringBranch string
	varStringImage  string
	varStringTZ     string

	// Cmd describes a docker command.
	Cmd = cobrax.NewCommand("docker", cobrax.WithRunE(dockerCommand))
)

func init() {
	dockerCmdFlags := Cmd.Flags()
	dockerCmdFlags.StringVarP(&varServiceName, "service_name", "s")
	dockerCmdFlags.StringVarPWithDefaultValue(&varServiceType, "service_type", "t", "rpc")
	dockerCmdFlags.StringVarPWithDefaultValue(&varStringBase, "base", "a", "alpine:latest")
	dockerCmdFlags.IntVarP(&varIntPort, "port", "p")
	dockerCmdFlags.StringVarP(&varStringHome, "home", "m")
	dockerCmdFlags.StringVarP(&varStringRemote, "remote", "r")
	dockerCmdFlags.StringVarP(&varStringBranch, "branch", "b")
	dockerCmdFlags.StringVarPWithDefaultValue(&varStringImage, "image", "i", "golang:1.20.3-alpine3.17")
	dockerCmdFlags.StringVarPWithDefaultValue(&varStringTZ, "tz", "z", "Asia/Shanghai")
}
