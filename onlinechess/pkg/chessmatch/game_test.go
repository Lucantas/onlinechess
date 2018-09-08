package chessmatch

import "testing"

func TestSocketRead(t *testing.T) {

}

func TestNewPlayer(t *testing.T) {
	client1 := &Client{}
	client2 := &Client{Match: "nein"}
	p1 := newPlayer(client1, "match_name")
	if p1.conn != client1 {
		t.Error("Expected p1.conn to be an empty client")
	}
	if p1.match != "match_name" {
		t.Error("Expected p1.match to be 'match_name'")
	}
	p2 := newPlayer(client2, "")
	if p2.conn != client2 {
		t.Error("Expected p2.conn to be an empty client")
	}
	if p2.match != "" {
		t.Error("Expected p2.match to be empty")
	}
}
