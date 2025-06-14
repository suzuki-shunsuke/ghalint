package token

import (
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/github"
	ghtoken "github.com/suzuki-shunsuke/urfave-cli-v3-util/keyring/ghtoken/cli"
	"github.com/urfave/cli/v3"
)

func New(logE *logrus.Entry) *cli.Command {
	return ghtoken.New(logE, github.KeyService)
}
