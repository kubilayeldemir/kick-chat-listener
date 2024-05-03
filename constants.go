package main

var ChannelAndChatIdMap = map[string]int{
	//"jahrein": 25314085,
	"naru":        25540277,
	"rraenee":     25951243,
	"chips":       25594578,
	"tolunayoren": 357571,
	"kaanflix":    7437325,
}

var ChatIdAndChannelMap = map[int]string{}

func init() {
	for k, v := range ChannelAndChatIdMap {
		ChatIdAndChannelMap[v] = k
	}
}
