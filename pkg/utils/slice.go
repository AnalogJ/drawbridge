package utils

import "sort"

func StringInSlice(slice []string, item string) bool {
	for _, val := range slice {
		if val == item {
			return true
		}
	}
	return false
}

func MapKeys(m map[string]interface{}) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

//func indexOf(answers []interface{}, item interface{}) (int) {
//	for k, v := range answers {
//		if v == item {
//			return k
//		}
//	}
//	return -1
//}
