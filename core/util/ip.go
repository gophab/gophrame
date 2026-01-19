package util

import "strings"

func GetRemoteHost(str string) string {
	segs := strings.Split(str, ":")
	if len(segs) > 1 {
		return strings.Join(segs[:len(segs)-1], ":")
	}
	return str
}

func GetRemotePort(str string) string {
	segs := strings.Split(str, ":")
	if len(segs) > 1 {
		return strings.Join(segs[len(segs)-1:], ":")
	}
	return str
}
