import os
from dotenv import load_dotenv
from flask import Flask, render_template, request, abort
from twilio.jwt.access_token import AccessToken
from twilio.jwt.access_token.grants import VideoGrant, ChatGrant
from twilio.rest import Client
from twilio.base.exceptions import TwilioRestException

env_file = os.environ.get("ENV_FILE")
if env_file:
    load_dotenv(env_file)

twilio_account_sid = os.environ.get("TWILIO_ACCOUNT_SID")
twilio_api_key_sid = os.environ.get("TWILIO_API_KEY_SID")
twilio_api_key_secret = os.environ.get("TWILIO_API_KEY_SECRET")
twilio_client = Client(twilio_api_key_sid, twilio_api_key_secret, twilio_account_sid)

frontend_env = {
    "MQTT_HOST": os.environ.get("MQTT_HOST", "localhost"),
    "MQTT_PORT": os.environ.get("MQTT_PORT", "1882"),
}

app = Flask(__name__)


def get_chatroom(name):
    for conversation in twilio_client.conversations.conversations.stream():
        if conversation.friendly_name == name:
            return conversation

    # a conversation with the given name does not exist ==> create a new one
    return twilio_client.conversations.conversations.create(friendly_name=name)


@app.route("/")
def index():
    return render_template("index.html", env=frontend_env)


@app.route("/office")
def office():
    return render_template("office.html", env=frontend_env)


@app.route("/login", methods=["POST"])
def login():
    username = request.get_json(force=True).get("username")
    if not username:
        abort(401)

    conversation = get_chatroom("My Room")
    try:
        conversation.participants.create(identity=username)
    except TwilioRestException as exc:
        # do not error if the user is already in the conversation
        if exc.status != 409:
            raise

    token = AccessToken(
        twilio_account_sid, twilio_api_key_sid, twilio_api_key_secret, identity=username
    )
    token.add_grant(VideoGrant(room="My Room"))
    token.add_grant(ChatGrant(service_sid=conversation.chat_service_sid))

    return {"token": token.to_jwt(), "conversation_sid": conversation.sid}


if __name__ == "__main__":
    host = "0.0.0.0"

    # Runs server with TLS if in a production environment,
    # otherwise runs on port specified in PORT environment variable (defaulting to 3000).
    environment = os.environ.get("ENV", "production")
    if environment == "production":
        app.run(host=host, port=443, ssl_context=("tls-cert.pem", "tls-key.pem"))
    else:
        port = int(os.environ.get("PORT", 3000))
        app.run(host=host, port=port)
