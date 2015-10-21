package realip

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/YunxiangHuang/vulcand-plugin/ipv4"
	"github.com/YunxiangHuang/vulcand-plugin/utils"
	"github.com/mailgun/vulcand/Godeps/_workspace/src/github.com/mailgun/log"
)

const (
	recursiveON  = "ON"
	recursiveOFF = "OFF"

	headerXFF = "X-FORWARDED-FOR"
	headerRIP = "REAL_IP"
	headerRAD = "REMOTE_ADDR"

	defaultAimHeader = "REALIP_XFF"
)

type RealIPHandler struct {
	cfg  RealIPMiddleware
	next http.Handler

	Header    string
	Recursive bool
	Whitelist ipv4.IPv4Segaments
	AimHeader string
}

type RealIPMiddleware struct {
	Recursive string
	Header    string
	Whitelist string
	Name      string
}

func New(re, he, wh, na string) (*RealIPMiddleware, error) {
	re = strings.TrimSpace(strings.ToUpper(re))
	he = strings.TrimSpace(strings.ToUpper(he))

	if re != recursiveON && re != recursiveOFF {
		return &RealIPMiddleware{}, fmt.Errorf("Config error - recursive: ", re)
	}

	if he != headerXFF && he != headerRAD && he != headerRIP && he != "" {
		return &RealIPMiddleware{}, fmt.Errorf("Config error - header: ", he)
	}

	if na == "" {
		na = defaultAimHeader
	}

	res := RealIPMiddleware{
		Recursive: re,
		Header:    he,
		Whitelist: wh,
		Name:      na,
	}

	return &res, nil
}

func (rih *RealIPHandler) setXForwardedFor(r *http.Request) {
	// rewrite NOT Append
	r.Header.Set(headerXFF, r.Header.Get(headerRAD))
}

func (rih *RealIPHandler) setAimHeaderWithXForwardedFor(aim string, r *http.Request) {
	xff := r.Header.Get(headerXFF)
	if xff != "" {
		list := utils.SplitWithoutSpace(xff, ",")
		if rih.Recursive {
			flag := true
			for i := len(list) - 1; i >= 0; i-- {
				tmpIP, err := ipv4.NewIPv4AddrFromString(list[i])
				if err != nil {
					continue
				}
				if !rih.Whitelist.IsInclude(tmpIP) {
					r.Header.Set(aim, tmpIP.String())
					flag = false
					break
				}
			}
			if flag {
				r.Header.Set(aim, list[0])
			}
		} else {
			r.Header.Set(aim, list[len(list)-1])
		}
	}
}

func (rim *RealIPMiddleware) NewHandler(next http.Handler) (http.Handler, error) {
	var res RealIPHandler
	if rim.Recursive == recursiveON {
		res.Recursive = true
	}
	if rim.Header == "" || rim.Header == headerRIP {
		res.Header = headerRAD
	} else {
		res.Header = rim.Header
	}

	if rim.Name == "" {
		res.AimHeader = defaultAimHeader
	} else {
		res.AimHeader = rim.Name
	}

	wList := utils.SplitWithoutSpace(rim.Whitelist, ",")
	for i := range wList {
		tmp, err := ipv4.NewIPv4SegamentFromString(wList[i])
		if err != nil {
			return &RealIPHandler{}, err
		}
		res.Whitelist = append(res.Whitelist, tmp)
	}

	res.Whitelist.Merge()
	res.next = next
	res.cfg = *rim
	return &res, nil
}

func (rih *RealIPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmp, _, _ := net.SplitHostPort(r.RemoteAddr)
	if reqIP, err := ipv4.NewIPv4AddrFromString(tmp); err == nil {
		if rih.Whitelist.IsInclude(reqIP) {
			switch rih.Header {
			case headerXFF:
				rih.setAimHeaderWithXForwardedFor(rih.AimHeader, r)
			default:
				r.Header.Set(rih.AimHeader, reqIP.String())
			}
		}
	} else {
		log.Errorf("Illegal Request", r)
	}
	rih.next.ServeHTTP(w, r)
}
