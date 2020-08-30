package lbry

import "testing"

func TestGetChannel(t *testing.T) {
	channel, err := GetChannelClaim("c9da929d12afe6066acc89eb044b552f0d63782a")
	if err != nil {
		t.Error(err)
	}
	if channel.Name != "@timcast" {
		t.Errorf("expected @canthareluscibarius, got %s", channel.ChannelName)
	}
}
