# CoffeeTalk

CoffeeTalk is a video chat application using WebRTC, WebSockets and MQTT. It is served as a web app, allowing users to stream video and audio peer-to-peer, and even play quizzes together!

**Table of Contents**

- [Project Structure](#project-structure)
- [Local Setup](#local-setup)
- [Credits](#credits)

## Project Structure

The application consists of 3 servers, all written in [Go](https://go.dev/), with a Docker container configured for each.

- `web/` contains a web server, serving the CoffeeTalk web app.
  - `server/` sets up the web server, defines the app's routes, and configures TLS and app environment.
  - `templates/` contains the HTML template files for the app's pages.
  - `static/` contains the static JavaScript and CSS files for the web app.
    - `scripts/` contains the client application logic, written in vanilla JavaScript and using JSDoc/TypeScript for type hinting. `main.js` is the app's entry point.
      - `_types/` contains TypeScript declarations for the application's custom types, to improve type hinting.
  - When served, the web app can be accessed through two routes:
    - `/` is the default app for users accessing CoffeeTalk from home, letting them choose a name and join the talk.
    - `/office` is intended for offices to set up CoffeeTalk on a computer in a break room, and have it automatically join the talk when motion is detected (using the [Diffy.js](https://github.com/maniart/diffyjs#readme) library).
- `webrtc/` contains a WebRTC signaling server for coordinating peer-to-peer video and audio streaming between clients. It establishes WebSocket connections with clients for persistent two-way communication and message forwarding.
  - `signaling/` sets up the socket connection API and defines the signaling messages passed between clients.
- `mqtt/` contains a server for running quiz sessions over MQTT.
  - `broker/` wraps around the [mochi-co/mqtt](https://github.com/mochi-co/mqtt#readme) package to set up an MQTT broker.
  - `quiz/` defines a state machine for running quiz sessions, publishing questions and answers to the broker.
- `stm/` contains a Go package with utility types and functions for setting up state machines. The documentation can be read in `stm.go`, or on [pkg.go.dev](https://pkg.go.dev/github.com/dcs-team4/coffeetalk/stm).

The project uses Docker Compose to coordinate containers, with a config for local development defined in `docker-compose.yml`, and a production config in `docker-compose-prod.yml`. The system has been deployed on a [DigitalOcean](https://www.digitalocean.com/) Virtual Private Server, but could be deployed anywhere that supports Docker.

For production, the servers expect a TLS certificate (`tls-cert.pem`) and key (`tls-key.pem`) in their respective `tls` directories (`web/server/tls`, `webrtc/signaling/tls`, `mqtt/broker/tls`). This is required, since the web browser can only access the webcam when the app is served over HTTPS.

The below deployment diagram shows the components in the system and the relations between them.

![deployment-diagram](https://raw.githubusercontent.com/dcs-team4/coffeetalk/docs/assets/deployment-diagram.png)

## Local setup

To run the project locally, follow these steps:

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

### Type Hinting

The web app uses [JSDoc](https://www.typescriptlang.org/docs/handbook/jsdoc-supported-types.html) for type hinting its JavaScript files through comments. It also uses TypeScript for custom type declarations, configured in `web/tsconfig.json`. To get full [type hinting of external libraries](https://code.visualstudio.com/docs/nodejs/working-with-javascript#_typings-and-automatic-type-acquisition) when developing in VSCode, install Node.js: https://nodejs.org/en/.

## Credits

Special thanks to:

- Frank Kraemer, for running the [Design of Communicating Systems](https://www.ntnu.edu/studies/courses/TTM4115) course at NTNU, which this project was a part of.
- Mozilla, for their great [documentation and tutorials for WebRTC](https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API/Signaling_and_video_calling), which was invaluable to our implementation.
- Gorilla, for the [Go WebSocket package](https://github.com/gorilla/websocket#readme).
- `mochi-co`, for the [MQTT broker package](https://github.com/mochi-co/mqtt#readme).
- The Eclipse Foundation, for the [Paho JavaScript MQTT client](https://www.eclipse.org/paho/index.php?page=clients/js/index.php).
- `maniart`, for the [Diffy.js library](https://github.com/maniart/diffyjs#readme).
