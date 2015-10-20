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
	SetProxyHeader string
}

func New(sp string) (*AddHeaderMiddleware, error) {
	hs := utils.SplitWithoutSpace(sp, ",")
	if sp == "" || len(hs) <= 0 {
		return &AddHeaderMiddleware{}, fmt.Errorf("Must set at least 1 Key-Value pair.")
	}
	return &AddHeaderMiddleware{
		SetProxyHeader: sp,
	}, nil
}

func (ahm *AddHeaderMiddleware) NewHandler(next http.Handler) (http.Handler, error) {
	var res AddHeaderHandler
	res.SetProxy = make(map[string]string)
	res.next = next
	hs := utils.SplitWithoutSpace(ahm.SetProxyHeader, ",")

	for i := range hs {
		tmp := utils.SplitWithoutSpace(hs[i], ":")
		if len(tmp) != 2 || tmp[1] == "" {
			return &res, fmt.Errorf("Format error: ", hs[i])
		} else {
			res.SetProxy[tmp[0]] = tmp[1]
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
