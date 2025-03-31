package main

import (
	"bufio"
	"bytes"
	"net"
	"testing"
)

type mockConn struct {
	net.Conn
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func setup() {
	peers = make(map[string]*Peer)
}

func TestProcessCommandBalance(t *testing.T) {
	setup()
	peer := &Peer{id: "peer1", balance: 100, conn: &mockConn{}}
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	processCommand(peer, "balance", writer)
	writer.Flush()

	expected := "Balance: 100\n"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestProcessCommandPay(t *testing.T) {
	setup()
	peer1 := &Peer{id: "peer1", balance: 100, conn: &mockConn{}}
	peer2 := &Peer{id: "peer2", balance: 50, conn: &mockConn{}}
	peers["peer1"] = peer1
	peers["peer2"] = peer2

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	processCommand(peer1, "Pay peer2 30", writer)
	writer.Flush()

	if peer1.balance != 70 {
		t.Errorf("expected peer1 balance to be 70, got %d", peer1.balance)
	}
	if peer2.balance != 80 {
		t.Errorf("expected peer2 balance to be 80, got %d", peer2.balance)
	}
}

func TestTransfer(t *testing.T) {
	setup()
	peer1 := &Peer{id: "peer1", balance: 100, conn: &mockConn{}}
	peer2 := &Peer{id: "peer2", balance: 50, conn: &mockConn{}}
	peers["peer1"] = peer1
	peers["peer2"] = peer2

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	transfer(peer1, "peer2", 30, writer)
	writer.Flush()

	if peer1.balance != 70 {
		t.Errorf("expected peer1 balance to be 70, got %d", peer1.balance)
	}
	if peer2.balance != 80 {
		t.Errorf("expected peer2 balance to be 80, got %d", peer2.balance)
	}
}

func TestTransferTargetNotFound(t *testing.T) {
	setup()
	peer1 := &Peer{id: "peer1", balance: 100, conn: &mockConn{}}
	peers["peer1"] = peer1

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	transfer(peer1, "peer3", 30, writer)
	writer.Flush()

	expected := "Target peer not found\n"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}
