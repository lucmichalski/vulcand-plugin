package addheader

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/YunxiangHuang/vulcand-plugin/utils"
)

const (
	HeaderFlag = "$"
)

type AddHeaderHandler struct {
	cfg  AddHeaderMiddleware
	next http.Handler

	SetProxy map[string]string
}

type AddHeaderMiddleware struct {
	SetProxyHeader []string
}

func New(sp []string) (*AddHeaderMiddleware, error) {
	if len(sp) <= 0 {
		return &AddHeaderMiddleware{}, fmt.Errorf("Must set at least 1 Key-Value pair :", sp)
	}
	return &AddHeaderMiddleware{
		SetProxyHeader: sp,
	}, nil
}

func (ahm *AddHeaderMiddleware) NewHandler(next http.Handler) (http.Handler, error) {
	var res AddHeaderHandler
	res.SetProxy = make(map[string]string)
	res.next = next

	for i := range ahm.SetProxyHeader {
		tmp := utils.SplitWithoutSpace(ahm.SetProxyHeader[i], ":")
		switch len(tmp) {
		case 2:
			res.SetProxy[strings.ToUpper(tmp[0])] = tmp[1]
		case 1:
			res.SetProxy[strings.ToUpper(tmp[0])] = ""
		default:
			return &res, fmt.Errorf("Format error: ", tmp)
		}
	}

	return &res, nil
}

func (ahh *AddHeaderHandler) SetProxyHeader(r *http.Request) {
	for k, v := range ahh.SetProxy {
		if strings.HasPrefix(v, HeaderFlag) {
			tmp := r.Header.Get(strings.TrimPrefix(v, HeaderFlag))
			r.Header.Set(k, tmp)
		} else {
			r.Header.Set(k, v)
		}
		if k == "X-FORWARDED-FOR" {
			r.Header.Set("REALIP", "AH_"+(r.Header.Get("REALIP")))
		}
	}
}

func (ahh *AddHeaderHandler) SetResponseHeader(w http.ResponseWriter) {
	// TODO: now do nothing
}

func (ahh *AddHeaderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ahh.SetProxyHeader(r)
	ahh.next.ServeHTTP(w, r)
	ahh.SetResponseHeader(w)
}
