package scatter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	_ "strconv"
	"time"
)

type BaseShuffle struct {
}

func (this BaseShuffle) RunShuffle(request *data_type.Request, IsOk func(data_type.Request, int) (int, bool)) int {
	defer common.Timer("ScatterMain", &(request.Timer_log), time.Now())
	for idx, _ := range request.Livelist {
		scope, is_ok := IsOk(*request, idx)
		if is_ok {
			continue
		}
		//确定需要交换位置
		change_idx := idx + scope
		for i := change_idx; i < len(request.Livelist); i++ {
			//需要交换位置满足条件，遍历查找
			request.Livelist[idx], request.Livelist[i] = request.Livelist[i], request.Livelist[idx]
			if _, is_ok := IsOk(*request, idx); is_ok {
				break
			} else {
				//如果失败，交换回来
				request.Livelist[idx], request.Livelist[i] = request.Livelist[i], request.Livelist[idx]
			}
		}
		if _, is_ok := IsOk(*request, idx); !is_ok {
			//交换完后仍不满足，打散结束
			break
		}
	}
	logs.Flush()
	return 0
}
