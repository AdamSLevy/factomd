package controlPanel_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/FactomProject/factomd/controlPanel"
	"github.com/FactomProject/factomd/p2p"
)

var _ = fmt.Sprintf("")

func TestFormatDuration(t *testing.T) {
	initial := time.Now().Add(-49 * time.Hour)
	if FormatDuration(initial) != "2 days" {
		t.Errorf("Time display incorrect : days")
	}

	initial = time.Now().Add(-25 * time.Hour)
	if FormatDuration(initial) != "1 day" {
		t.Errorf("Time display incorrect : days")
	}

	initial = time.Now().Add(-23 * time.Hour)
	if FormatDuration(initial) != "23 hrs" {
		t.Errorf("Time display incorrect : hrs")
	}

	initial = time.Now().Add(-1 * time.Hour)
	if FormatDuration(initial) != "1 hr" {
		t.Errorf("Time display incorrect : hr")
	}

	initial = time.Now().Add(-59 * time.Minute)
	if FormatDuration(initial) != "59 mins" {
		t.Errorf("Time display incorrect : mins")
	}

	initial = time.Now().Add(-1 * time.Minute)
	if FormatDuration(initial) != "1 min" {
		t.Errorf("Time display incorrect : min")
	}

	initial = time.Now().Add(-30 * time.Second)
	if FormatDuration(initial) != "30 secs" {
		t.Errorf("Time display incorrect : secs")
	}
}

func TestTallyTotals(t *testing.T) {
	cm := NewConnectionsMap()
	var i uint32
	for i = 0; i < 10; i++ {
		cm.Connect(fmt.Sprintf("%d", i), NewP2PConnection(i, i, i, i, fmt.Sprintf("%d", i), i))
	}
	for i = 10; i < 20; i++ {
		cm.Disconnect(fmt.Sprintf("%d", i), NewP2PConnection(i, i, i, i, fmt.Sprintf("%d", i), i))
	}
	cm.TallyTotals()
	if cm.Totals.BytesSentTotal != 190 {
		t.Errorf("Byte Sent does not match")
	}
	if cm.Totals.BytesReceivedTotal != 190 {
		t.Errorf("Byte Received does not match")
	}
	if cm.Totals.MessagesSent != 190 {
		t.Errorf("Msg Sent does not match")
	}
	if cm.Totals.MessagesReceived != 190 {
		t.Errorf("Msg Received does not match")
	}
	if cm.Totals.PeerQualityAvg != 4 {
		t.Errorf("Peer Quality does not match %d", cm.Totals.PeerQualityAvg)
	}

	for key := range cm.GetConnectedCopy() {
		cm.RemoveConnection(key)
	}
	for key := range cm.GetDisconnectedCopy() {
		cm.RemoveConnection(key)
	}
	cm.TallyTotals()
	if cm.Totals.BytesSentTotal != 0 {
		t.Errorf("Byte Sent does not match")
	}
	if cm.Totals.BytesReceivedTotal != 0 {
		t.Errorf("Byte Received does not match")
	}
	if cm.Totals.MessagesSent != 0 {
		t.Errorf("Msg Sent does not match")
	}
	if cm.Totals.MessagesReceived != 0 {
		t.Errorf("Msg Received does not match")
	}
	if cm.Totals.PeerQualityAvg != 0 {
		t.Errorf("Peer Quality does not match %d", cm.Totals.PeerQualityAvg)
	}
}

func NewP2PConnection(bs uint32, br uint32, ms uint32, mr uint32, addr string, pq uint32) *p2p.ConnectionMetrics {
	pc := new(p2p.ConnectionMetrics)
	pc.MomentConnected = time.Now()
	pc.BytesSent = bs
	pc.BytesReceived = br
	pc.MessagesSent = ms
	pc.MessagesReceived = mr
	pc.PeerAddress = addr
	pc.PeerQuality = int32(pq)

	return pc
}
