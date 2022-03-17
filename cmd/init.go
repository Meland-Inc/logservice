package cmd

import (
	"github.com/urfave/cli/v2"
)

/// 这里对一些服务进行初始化聚合
func xinit(c *cli.Context) error {
	inits := []func() error{}

	for _, init := range inits {
		if err := init(); err != nil {
			return err
		}
	}

	return nil
}
