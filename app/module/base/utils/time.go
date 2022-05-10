package utils

import (
	"container/list"
	"sort"
	"time"
)

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ListToArray
// list对象转数组
func ListToArray(list *list.List) []interface{} {
	var listLen = list.Len()
	if listLen == 0 {
		return nil
	}
	var arr []interface{}
	for e := list.Front(); e != nil; e = e.Next() {
		arr = append(arr, e.Value)
	}
	return arr
}

// ExistsDuplicateInStringsArr 字符串数组中是否存在重复元素
func ExistsDuplicateInStringsArr(arr []string) bool {
	length := len(arr)
	sort.Strings(arr)
	for i := 1; i < length; i++ {
		if arr[i-1] == arr[i] {
			return true
		}
	}
	return false
}
