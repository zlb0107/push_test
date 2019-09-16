package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type FilmFilterHasShown struct {
}

func init() {
	var rp FilmFilterHasShown
	Filter_map["FilmFilterHasShown"] = rp
	logs.Warn("in filter FilmFilterHasShown init")
}
func (rp FilmFilterHasShown) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	if request.Film_has_shown_list == nil {
		return false
	}
	if _, is_in := request.Film_has_shown_list[info.Uid]; is_in {
		return true
	}
	return false
}
