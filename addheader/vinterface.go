package addheader

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mailgun/vulcand/plugin"
)

const Type = "addheader"

func GetSpec() *plugin.MiddlewareSpec {
	return &plugin.MiddlewareSpec{
		Type:      Type,
		FromOther: FromOther,
		FromCli:   FromCli,
		CliFlags:  CliFlags(),
	}
}

func FromOther(ahm AddHeaderMiddleware) (plugin.Middleware, error) {
	return New(ahm.SetProxyHeader)
}

func FromCli(c *cli.Context) (plugin.Middleware, error) {
	if !c.IsSet("setproxyheader") {
		if !c.IsSet("S") {
			return &AddHeaderMiddleware{}, fmt.Errorf("Miss Argument: setproxyheader.")
		} else {
			return New(c.StringSlice("S"))
		}
	}
	return New(c.StringSlice("setproxyheader"))
}

func CliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringSliceFlag{"setproxyheader, S", &cli.StringSlice{}, "set proxy header and value", ""},
	}
}
