package multi_exp

import (
	"bytes"
	logs "github.com/cihub/seelog"
)

func Get_bucket(uid string, seed uint) uint {
	final := MyMurmurhash2([]byte(uid), uint32(seed))
	bucket := uint(final) % 100
	return bucket
}

func Get_expid(uid string, bucket uint, level *LevelStruct) string {
	expid := ""
	priority := 0
	for _, exp := range level.Exps {
		switch exp.Generate_type {
		case UseStartEnd:
			{
				if bucket >= exp.Start && bucket <= exp.End {
					if exp.Priority > priority {
						expid = exp.Expid
						priority = exp.Priority
					}
				}
			}
		case UseUids, UseRedis:
			{
				if _, is_in := exp.Uids_map[uid]; is_in {
					if exp.Priority > priority {
						expid = exp.Expid
						priority = exp.Priority
					}
				}
			}
		default:
			{
				logs.Error("it is wired, type:", exp.Generate_type)
			}
		}
	}
	return expid
}
func Get_expids(uid string) string {
	var Expid bytes.Buffer
	for _, level := range Exp_handler.Levels {
		bucket := Get_bucket(uid, level.Seed)
		expid := Get_expid(uid, bucket, &level)
		if expid != "" {
			if Expid.Len() != 0 {
				Expid.WriteString(",")
			}
			Expid.WriteString(expid)
		}
	}
	return Expid.String()
}
