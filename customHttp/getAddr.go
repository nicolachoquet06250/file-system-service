package customHttp

import (
	"filesystem_service/flags"
	"fmt"
	"strconv"
	"strings"
)

func GetAddr() string {
	fl := flags.GetFlags()

	host := fl.GetHost()
	if host != "" && strings.Contains(host, ":") {
		host = fmt.Sprintf("[%s]", host)
	}

	return fmt.Sprintf("%v:%v", host, strconv.Itoa(fl.GetPort()))
}
