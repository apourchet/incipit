package utils

import "github.com/parnurzeal/gorequest"

func GetRpc(url string, data interface{}, endStruct interface{}) error {
	_, _, errs := gorequest.New().Get(url).
		Send(data).EndStruct(endStruct)
	if len(errs) != 0 {
		return errs[0]
	}
	return nil
}
