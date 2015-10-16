package realip

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/YunxiangHuang/ipv4"
)

const (
	recursiveON  = "ON"
	recursiveOFF = "OFF"

	headerXFF = "X-FORWARDED-FOR"
	headerRIP = "REAL_IP"
	headerRAD = "REMOTE_ADDR"
)

type RealIPHandler struct {
	cfg  RealIPMiddleware
	next http.Handler

	Header    string
	Recursive bool
	Whitelist ipv4.IPv4Segaments
}

type RealIPMiddleware struct {
	Recursive string
	Header    string
	Whitelist string
}

func New(re, he, wh string) (*RealIPMiddleware, error) {
	re = strings.TrimSpace(strings.ToUpper(re))
	he = strings.TrimSpace(strings.ToUpper(he))

	if re != recursiveON && re != recursiveOFF {
		return &RealIPMiddleware{}, fmt.Errorf("Config error - recursive: ", re)
	}

	if he != headerXFF && he != headerRAD && he != headerRIP && he != "" {
		return &RealIPMiddleware{}, fmt.Errorf("Config error - header: ", he)
	}

	res := RealIPMiddleware{
		Recursive: re,
		Header:    he,
		Whitelist: wh,
	}

	fmt.Println("New res: \n", res, "\n\n")

	return &res, nil
}

func (rih *RealIPHandler) setXForwardedFor(r *http.Request) {
	// rewrite NOT Append
	fmt.Printf("bSet XFF:%v\tRAD:%v\n", r.Header.Get(headerXFF), r.Header.Get(headerRAD))
	r.Header.Set(headerXFF, r.Header.Get(headerRAD))
	fmt.Printf("aSet XFF:%v\tRAD:%v\n", r.Header.Get(headerXFF), r.Header.Get(headerRAD))
}

func (rih *RealIPHandler) setRemoteAddrWithXForwardedFor(r *http.Request) {
	fmt.Println("*******")
	xff := r.Header.Get(headerXFF)
	if xff != "" {
		list := strings.Split(xff, ",")
		if rih.Recursive {
			flag := true
			for i := len(list) - 1; i >= 0; i-- {
				tmpIP, err := ipv4.NewIPv4AddrFromString(list[i])
				if err != nil {
					continue
				}
				if !rih.Whitelist.IsInclude(tmpIP) {
					r.Header.Set(headerRAD, tmpIP.String())
					flag = false
					break
				}
			}
			if flag {
				r.Header.Set(headerRAD, list[0])
			}
		} else {
			r.Header.Set(headerRAD, list[len(list)-1])
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

	wList := strings.Split(rim.Whitelist, ",")
	for i := range wList {
		tmp, err := ipv4.NewIPv4SegamentFromString(wList[i])
		if err != nil {
			return &RealIPHandler{}, err
		}
		res.Whitelist = append(res.Whitelist, tmp)
	}
	res.Whitelist = ipv4.IPv4SegamentsMerge(res.Whitelist)

	fmt.Printf("New RealIP: %v\n", res)

	res.next = next
	res.cfg = *rim

	fmt.Println("\n\nNew Handler: ", *rim)

	return &res, nil
}

func (rih *RealIPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("============")
	fmt.Printf("RAD:%v\tXFF:%v\n", r.Header.Get(headerRAD), r.Header.Get(headerXFF))
	reqIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	fmt.Printf("### RealIP: %v\n", rih.cfg)
	switch rih.Header {
	case headerXFF:
		rih.setRemoteAddrWithXForwardedFor(r)
	default:
		r.Header.Set(headerRAD, reqIP)
	}
	rih.setXForwardedFor(r)
	fmt.Printf("RAD:%v\tXFF:%v\n", r.Header.Get(headerRAD), r.Header.Get(headerXFF))

	rih.next.ServeHTTP(w, r)
}
