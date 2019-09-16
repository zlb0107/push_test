package redis_queue

type HasShownTantan struct {
	RedisQueue
}

var HasShownTantanHandler HasShownTantan

func init() {
	HasShownTantanHandler = HasShownTantan{RedisQueue{queueLen: 500, keyPrefix: "tantan_shown_", chanLen: 10000}}
	HasShownTantanHandler.Info_chan = make(chan ZidInfoStruct, HasShownTantanHandler.chanLen)
	go HasShownTantanHandler.Deal_chan()
}
