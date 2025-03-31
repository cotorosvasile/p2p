
```markdown
# Peer-to-Peer Payment System

This is a simple peer-to-peer payment system implemented in Go. It allows peers to connect to a server, check their balance, and transfer funds to other peers.

## Features

- **Server**: Listens for incoming peer connections.
- **Client**: Connects to the server and allows the user to interact with the system.
- **Commands**:
  - `balance`: Check the current balance.
  - `Pay <peer_id> <amount>`: Transfer funds to another peer.

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/cotorosvasile/p2p.git
   cd p2p
   ```

2. Build the application:
   ```sh
   go build -o start-peer
   ```

## Usage

### Starting the Server

To start the server, run:
```sh
./start-peer server <port>
```
Replace `<port>` with the port number you want the server to listen on.

### Starting a Client

To start a client, run:
```sh
./start-peer client <address> <peer_id>
```
Replace `<address>` with the server address (e.g., `localhost:8080`) and `<peer_id>` with a unique identifier for the peer.

### Example

1. Start the server:
   ```sh
   ./start-peer server 8080
   ```

2. Start a client:
   ```sh
   ./start-peer client localhost:8080 peer1
   ```

3. In the client, you can now use the following commands:
    - `balance`: Check your balance.
    - `Pay peer2 50`: Transfer 50 units to `peer2`.

## Running Tests

To run the unit tests, use the following command:
```sh
go test ./...
```
