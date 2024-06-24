# WebRTC-SFU Server
This repository contains a Go implementation of a Selective Forwarding Unit (SFU) for WebRTC. The server uses WebSocket for signaling and the Gin framework for HTTP handling. It is designed to facilitate real-time communication applications by managing the connections and data transfer between peers efficiently.

# Features
- **WebRTC Integration**: Supports WebRTC peer connections for video and audio streams.
- **WebSocket Signaling**: Utilizes WebSocket for real-time signaling between clients and the server.
- **Gin Framework**: Uses Gin for handling HTTP requests and serving static files.
- **Selective Forwarding**: Efficiently manages and forwards media streams to connected peers.
- **ICE Candidate Management**: Handles ICE candidates to establish peer-to-peer connections.
- **Track Management**: Dynamically adds and removes media tracks and manages renegotiation.

### Getting Started

#### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/MajorTom3K1M/go-webrtc-sfu.git
   cd go-webrtc-sfu
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

#### Running the Server

1. Build and run the server:
   ```sh
   go run .
   ```

2. Open your browser and navigate to `http://localhost:8080`.

### Acknowledgements

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Pion WebRTC](https://github.com/pion/webrtc)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)
