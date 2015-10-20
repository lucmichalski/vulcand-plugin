package utils

import(
	"strings"
)

func SplitWithoutSpace(str string, flag string) []string {
		res := strings.Split(str, flag)
	for i := range res {
		res[i] = strings.TrimSpace(res[i])
	}
	return res
}