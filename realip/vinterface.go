package realip

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/eleme/vulcand/plugin"
)

const Type = "realip"

func GetSpec() *plugin.MiddlewareSpec {
	return &plugin.MiddlewareSpec{
		Type:      Type,
		FromOther: FromOther,
		FromCli:   FromCli,
		CliFlags:  CliFlags(),
	}
}

func FromOther(r RealIPMiddleware) (plugin.Middleware, error) {
	return New(r.Recursive, r.Header, r.Whitelist, r.Name)
}

func FromCli(c *cli.Context) (plugin.Middleware, error) {
	if !c.IsSet("recursive") || !c.IsSet("whitelist") {
		return &RealIPMiddleware{}, fmt.Errorf("Missing Arguments: recursive or whitelist.")
	}
	
	return New(c.String("recursive"), c.String("header"), c.String("whitelist"), "REALIP")
}

func CliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{"recursive, R", "", "Enable Recursive [ON|OFF]", ""},
		cli.StringFlag{"header, H", "REMOTE_ADDR", "Which set to X-FORWARDED-FOR [REMOTE_ADDR|X-FORWARDED-FOR]", ""},
		cli.StringFlag{"whitelist, W", "", "Whitelist, format: 1.1.1.1/24 or 1.1.1.1", ""},
		//cli.StringFlag{"name, N", "REALIP_XFF", "Realip will set this Key-Value to Header and not set X-FORWARDED-FOR", ""},
	}
}