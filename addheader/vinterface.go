package addheader

import (
	"fmt"

	"github.com/mailgun/vulcand/Godeps/_workspace/src/github.com/codegangsta/cli"
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
		return &AddHeaderMiddleware{}, fmt.Errorf("Missing Argument: setproxyheader.")
	}
	return New(c.String("setproxyheader"))
}

func CliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{"setproxyheader, S", "", "set proxy header and value like: key1:value1, key2:" + HeaderFlag + "X-FORWARDED-FOR", ""},
	}
}
