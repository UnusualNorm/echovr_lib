package messages

import (
	echovr "github.com/unusualnorm/echovr_lib"
)

var SNSChannelInfoResponseSymbol uint64 = echovr.GenerateSymbol("SNSChannelInfoResponse")

type Channel struct {
	ChannelUUID  string `json:"channeluuid"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Rules        string `json:"rules"`
	RulesVersion uint64 `json:"rules_version"`
	Link         string `json:"link"`
	Priority     uint64 `json:"priority"`
	RAD          bool   `json:"_rad"`
}

type ChannelInfo struct {
	Groups []Channel `json:"groups"`
}

type SNSChannelInfoResponse struct {
	ChannelInfo ChannelInfo
}

func (m *SNSChannelInfoResponse) Symbol() uint64 {
	return SNSChannelInfoResponseSymbol
}

func (m *SNSChannelInfoResponse) Stream(s *echovr.EasyStream) error {
	return s.StreamZlibEasyStream(func(decompressedS *echovr.EasyStream) error {
		return decompressedS.StreamJson(&m.ChannelInfo)
	})
}
