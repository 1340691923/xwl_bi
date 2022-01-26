package util

import jsoniter "github.com/json-iterator/go"

func Model2Map(m interface{}, needZeroByInt, needZeroByString bool) (res map[string]interface{}) {
	res = map[string]interface{}{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &res)
	for k, v := range res {
		switch v.(type) {
		case float64:
			if v.(float64) == 0 && !needZeroByInt {
				delete(res, k)
			}
		case int64:
			if v.(int64) == 0 && !needZeroByInt {
				delete(res, k)
			}
		case int:
			if v.(int) == 0 && !needZeroByInt {
				delete(res, k)
			}
		case string:
			if v.(string) == "" && !needZeroByString {
				delete(res, k)
			}
		}
	}
	return
}
