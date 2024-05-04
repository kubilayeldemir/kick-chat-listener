package main

var ChannelAndChatIdMap = map[string]int{
	"jahrein":     25314085,
	"naru":        25540277,
	"rraenee":     25951243,
	"chips":       25594578,
	"tolunayoren": 357571,
	"kaanflix":    7437325,
	"hype":        24495088,
	"elwind":      25240548,
	"elraenn":     25712360,
	"purplebixi":  25593921,
	"swaggybark":  25593949,
	"rammus53":    1292179,
	"uthenera":    7320801,
	"levo":        24906135,
	"eray":        10181332,
	"ebonivon":    25813693,
}

var ChatIdAndChannelMap = map[int]string{}

func init() {
	for k, v := range ChannelAndChatIdMap {
		ChatIdAndChannelMap[v] = k
	}
}
