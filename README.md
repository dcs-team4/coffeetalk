# CoffeeTalk

Inter-office communication system.

## Project Structure

The project consists of:

- A Python web server in `/frontend` hosting the static HTML/CSS/JavaScript files in `/frontend/static`.
  - This uses [Twilio](https://www.twilio.com/) for video streaming between clients.
- An MQTT broker and state machine for running quiz sessions.
  - `/mqtt/broker` wraps around the Go library [`mochi-co/mqtt`](https://github.com/mochi-co/mqtt)
    to set up an MQTT broker.
  - `/mqtt/quiz` contains a Go package that defines a state machine for quiz sessions, as well as the configuration of quiz questions. It uses the MQTT broker to publish questions to the appropriate topics.
- `/stm`, a Go package with generic utility types and functions for composing state machines.

## Local setup

1. Install and open Docker Desktop: https://www.docker.com/products/docker-desktop/

2. Clone the repository (from the terminal):

```
git clone https://github.com/dcs-team4/coffeetalk.git
```

3. Navigate to the cloned folder:

```
cd coffeetalk
```

4. Build and run the Docker containers:

```
docker compose up --build
```

If this fails, check that your Docker installation is configured correctly.

The web application should now be accessible at `localhost:3000`, with the MQTT broker served through WebSocket at `localhost:1882`, and TCP at `localhost:1883`. The video streaming is done through WebRTC, coordinated by a remote [Twilio](https://www.twilio.com/) server.
