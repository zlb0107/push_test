package merger

import (
	"go_common_lib/data_type"
)

func De_duplication(array_map map[string][]data_type.LiveInfo, merge_map map[string]MergeStruct, priority_list []string) {
	uid_map := make(map[string]int)
	for _, name := range priority_list {
		list, is_in := array_map[name]
		if !is_in {
			continue
		}
		ms, is_in := merge_map[name]
		if !is_in {
			continue
		}
		delete_uid := make([]int, 0)
		for i, info := range list {
			if _, is_in := uid_map[info.Uid]; is_in {
				delete_uid = append(delete_uid, i)
			} else {
				uid_map[info.Uid] = 1
			}
		}
		for i := len(delete_uid); i > 0; i-- {
			list = append(list[:i], list[i+1:]...)
		}
		ms.Length -= len(delete_uid)
		array_map[name] = list
		merge_map[name] = ms
	}
}
func Merge_token(array_map map[string][]data_type.LiveInfo) {
	//分两遍，第一遍找出所有token，第二部重新给token赋值
	uid_token_map := make(map[string]string)
	for _, list := range array_map {
		for _, info := range list {
			token, is_in := uid_token_map[info.Uid]
			if is_in {
				token += info.Token
			} else {
				token = info.Token
			}
			uid_token_map[info.Uid] = token
		}
	}
	//第二次遍历
	for name, list := range array_map {
		for idx, info := range list {
			list[idx].Token = uid_token_map[info.Uid]
		}
		array_map[name] = list
	}
}
