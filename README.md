# CoffeeTalk

Inter-office communication system using state machines.

## Setup

Clone the repo (from the terminal):

```
git clone https://github.com/dcs-team4/coffeetalk.git
```

Navigate to the right folder:

```
cd coffeetalk
```

Install Python dependencies:

```
pip install -r requirements.txt
```

## Project Structure

The project consists of:

- A Python web server in `/server` hosting the static HTML/CSS/JavaScript files in `/frontend`.
  - This uses [Twilio](https://www.twilio.com/) for video streaming between clients.
- An MQTT broker for running quiz sessions.
  - `/mqtt/broker` wraps around the Go library [`mochi-co/mqtt`](https://github.com/mochi-co/mqtt)
    to set up an MQTT broker.
  - `/mqtt/quiz` contains a Go package that defines a state machine for quiz sessions, as well as the configuration of quiz questions. It uses the MQTT broker to publish questions to the appropriate topics.
  - `/stm` contains a Go package with generic utility types for state machines.
