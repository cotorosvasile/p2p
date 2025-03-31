package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Peer struct {
	id      string
	conn    net.Conn
	balance int
}

var (
	peers = make(map[string]*Peer)
	mutex sync.Mutex
)

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)
	if !scanner.Scan() {
		conn.Close()
		return
	}
	peerID := strings.TrimSpace(scanner.Text())

	peer := &Peer{
		id:      peerID,
		conn:    conn,
		balance: 0,
	}

	mutex.Lock()
	peers[peerID] = peer
	mutex.Unlock()

	fmt.Println("Peer connected:", peerID)

	for scanner.Scan() {
		text := scanner.Text()
		processCommand(peer, text, writer)
		writer.Flush()
	}

	mutex.Lock()
	delete(peers, peer.id)
	mutex.Unlock()
	fmt.Println("Peer disconnected:", peerID)
}

func processCommand(peer *Peer, text string, writer *bufio.Writer) {
	parts := strings.Fields(text)
	if len(parts) < 1 {
		return
	}

	switch parts[0] {
	case "balance":
		fmt.Fprintf(writer, "Balance: %d\n", peer.balance)
	case "Pay":
		if len(parts) < 3 {
			fmt.Fprintf(writer, "Usage: Pay <peer_id> <amount>\n")
			return
		}
		targetID := parts[1]
		amount, err := strconv.Atoi(parts[2])
		if err != nil || amount <= 0 {
			fmt.Fprintf(writer, "Invalid amount\n")
			return
		}
		transfer(peer, targetID, amount, writer)
	default:
		fmt.Fprintf(writer, "Unknown command\n")
	}
}

func transfer(sender *Peer, targetID string, amount int, writer *bufio.Writer) {
	mutex.Lock()
	target, exists := peers[targetID]
	mutex.Unlock()

	if !exists {
		fmt.Fprintf(writer, "Target peer not found\n")
		return
	}

	sender.balance -= amount
	target.balance += amount

	senderWriter := bufio.NewWriter(sender.conn)
	targetWriter := bufio.NewWriter(target.conn)

	fmt.Fprintf(senderWriter, "Sent %d to %s\n", amount, targetID)
	fmt.Fprintf(targetWriter, "Received %d from %s\n", amount, sender.id)

	senderWriter.Flush()
	targetWriter.Flush()
}

func startPeer(address, id string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s\n", id)
	fmt.Println("Welcome to your peering relationship! Your ID:", id)

	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		serverScanner := bufio.NewScanner(conn)
		for serverScanner.Scan() {
			fmt.Println(serverScanner.Text())
		}
	}()

	for {
		fmt.Print("> ")
		if scanner.Scan() {
			text := scanner.Text()
			fmt.Fprintf(conn, "%s\n", text)
		}
	}
}

func startServer(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("Waiting for peer connections on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  Server: ./start-peer server <port>")
		fmt.Println("  Client: ./start-peer client <address> <peer_id>")
		return
	}

	role := os.Args[1]
	param := os.Args[2]

	if role == "server" {
		startServer(param)
	} else if role == "client" {
		if len(os.Args) < 4 {
			fmt.Println("Usage: ./start-peer client <address> <peer_id>")
			return
		}
		peerID := os.Args[3]
		startPeer(param, peerID)
	} else {
		fmt.Println("Invalid role, use 'server' or 'client'")
	}
}
