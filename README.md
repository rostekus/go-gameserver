# Game Server with Actor Model in Golang

## Introduction

This repository contains a game server implementation using the actor model paradigm in Golang. The server is designed to handle concurrent interactions between players, game objects, and other entities in a scalable and efficient manner.

## Actor Model Paradigm

The actor model is a concurrency model that treats actors as fundamental units of computation. Actors are independent entities that can communicate with each other by exchanging messages. In this game server, we utilize the actor model to manage game entities as actors and facilitate communication and interaction between them.

## Features

- Concurrent Gameplay: The actor model allows for concurrent execution of game actions, enabling seamless and responsive gameplay even under heavy loads.

- Scalability: With the actor model, the server can easily scale to handle an increasing number of players and entities without sacrificing performance.

- Fault Tolerance: Each actor is isolated and only communicates through messages, which enhances fault tolerance and system reliability.

- Flexibility: The actor model simplifies the process of adding new game features or modifying existing ones, making the server highly adaptable.

- Websocket Communication: Communication between clients and the server is established using WebSockets, providing real-time, bidirectional communication for multiplayer gaming experiences.

- Protocol Buffers (protobuf): To optimize data serialization and deserialization, Protocol Buffers are used for message exchange between clients and the server.

![Untitled Diagram drawio](https://github.com/rostekus/go-gameserver/assets/34031791/b61c8244-f5ac-4c53-892f-8f22b381d309)

### Prerequisites

- Golang installed on your system.

### Installation

1. Clone the repository: `git clone https://github.com/your-username/game-server.git`

2. Change directory: `cd game-server`

3. Build the server: `go build ./cmd/server`

4. Build the client: `go build ./cmd/client`

### Usage

1. Run the server: `./server`

2. Run the client: `./client`

## License

This project is licensed under the [MIT License](LICENSE).

