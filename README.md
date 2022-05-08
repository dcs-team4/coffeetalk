# CoffeeTalk

CoffeeTalk is an inter-office communication system using WebRTC, WebSockets and MQTT. It is served as a web app, allowing users to stream video and audio peer-to-peer, and even play quizzes together!

**Table of Contents**

- [Project Structure](#project-structure)
- [Local Setup](#local-setup)
- [Credits](#credits)

## Project Structure

The project provides 3 servers, all written in [Go](https://go.dev/), with a Docker container configured for each.

- `/web` contains a web server, serving the CoffeeTalk web app from the `templates` and `static` folders. The web app itself is written in vanilla HTML/CSS/JavaScript, with most of its logic in `/web/static/scripts`.
  - `/server` sets up the web server, defines the app's routes, and configures TLS and web app environment.
  - When served, the web app can be accessed through two routes:
    - `/` is the default app for users accessing CoffeeTalk from home, letting them choose a name and join the talk.
    - `/office` is intended for offices to set up CoffeeTalk on a computer in a break room, and have it automatically join the talk when motion is detected (the [Diffy.js](https://github.com/maniart/diffyjs#readme) library is used for this).
- `/webrtc` contains a WebRTC signaling server for coordinating peer-to-peer video and audio streaming between clients. It establishes WebSocket connections with clients for persistent two-way communication and message forwarding.
  - `/signals` contains the main signaling logic.
- `/mqtt` contains a server for running quiz sessions over MQTT.
  - `/broker` wraps around the [mochi-co/mqtt](https://github.com/mochi-co/mqtt#readme) package to set up an MQTT broker.
  - `/quiz` defines a state machine for running quiz sessions, publishing questions and answers to the broker.
- `/stm` contains a Go package with utility types and functions for setting up state machines. The documentation can be read in `stm.go`, or on [pkg.go.dev](https://pkg.go.dev/github.com/dcs-team4/coffeetalk/stm).

The project uses Docker Compose to coordinate containers, with a config for local development defined in `docker-compose.yml`, and a production config in `docker-compose-prod.yml`. The system has been deployed on a [DigitalOcean](https://www.digitalocean.com/) Virtual Private Server, but could be deployed anywhere that supports Docker.

For production, the servers expect a TLS certificate (`tls-cert.pem`) and key (`tls-key.pem`) in their respective `tls` directories (`web/server/tls`, `webrtc/signals/tls`, `mqtt/broker/tls`). This is required, since the web browser can only access the webcam when the app is served over HTTPS.

The below deployment diagram shows the components in the system and the relations between them.

![deployment-diagram](https://raw.githubusercontent.com/dcs-team4/coffeetalk/docs/assets/deployment-diagram.png)

## Local setup

1. Install and open Docker Desktop: https://www.docker.com/products/docker-desktop/

2. Clone the repository, and navigate into it:

```
git clone https://github.com/dcs-team4/coffeetalk.git
cd coffeetalk
```

3. Build and run the Docker containers:

```
docker compose up --build
```

The web application should now be accessible at `localhost:3000`, with the WebRTC signaling server listening on `localhost:8000`, and the MQTT broker served over WebSocket at `localhost:1882` and TCP at `localhost:1883`.

## Credits

Special thanks to:

- Frank Kraemer, for running the [Design of Communicating Systems](https://www.ntnu.edu/studies/courses/TTM4115) course at NTNU, which this project was a part of.
- Mozilla, for their great [documentation and tutorials for WebRTC](https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API/Signaling_and_video_calling), which was invaluable to our implementation.
- `maniart`, for the [Diffy.js library](https://github.com/maniart/diffyjs#readme).
- Gorilla, for the [Go WebSocket package](https://github.com/gorilla/websocket#readme).
- `mochi-co`, for the [MQTT broker package](https://github.com/mochi-co/mqtt#readme).
- The Eclipse Foundation, for the [Paho JavaScript MQTT Client](https://www.eclipse.org/paho/index.php?page=clients/js/index.php).
