import { env } from "./env.js";
import { DOM } from "./dom.js";

/** Shared prefix for MQTT quiz topics. */
const MQTT_TOPIC_PREFIX = "coffeetalk/quiz";

/**
 * MQTT quiz topics to listen and post on.
 * Should be updated if the quiz server topic configuration changes.
 */
const mqttTopics = {
  QUESTIONS: `${MQTT_TOPIC_PREFIX}/questions`,
  ANSWERS: `${MQTT_TOPIC_PREFIX}/answers`,
  STATUS: `${MQTT_TOPIC_PREFIX}/status`,
};

/**
 * MQTT quiz messages that the server expects.
 * Should be updated if the quiz server message configuration changes.
 */
const mqttMessages = {
  START: "start-quiz",
  END: "end-quiz",
};

/**
 * Connection to the MQTT broker.
 * Undefined if uninitialized.
 * @type {Paho.MQTT.Client | undefined}
 */
let mqttClient;

/** Connects to the MQTT broker and sets up message listeners. */
export function connectMQTT() {
  // Passes "" as MQTT client ID, as that tells the server to generate a unique ID for us.
  mqttClient = new Paho.MQTT.Client(env.MQTT_HOST, parseInt(env.MQTT_PORT), "");
  mqttClient.onMessageArrived = handleMQTTMessage;
  mqttClient.onConnectionLost = ({ errorCode, errorMessage }) => {
    // Error code 0 means no error.
    if (errorCode !== 0) {
      console.log("Connection to MQTT broker lost:", errorMessage);
    }
  };

  mqttClient.connect({
    onSuccess: () => {
      console.log("Successfully connected to MQTT broker.");
      mqttClient?.subscribe(`${MQTT_TOPIC_PREFIX}/#`);
    },
    onFailure: ({ errorMessage }) => {
      console.log("Failed to connect to MQTT broker:", errorMessage);
    },
    useSSL: env.ENV === "production",
  });
}

/**
 * Handles the given MQTT message according to its topic.
 * @param {Paho.MQTT.Message} message
 */
function handleMQTTMessage(message) {
  console.log("MQTT message received:", {
    topic: message.destinationName,
    message: message.payloadString,
  });

  switch (message.destinationName) {
    case mqttTopics.QUESTIONS:
      DOM.quizQuestion().innerText = message.payloadString;
      DOM.quizAnswer().innerText = "";

      // Initializes quiz view on receiving question rather than receiving start message,
      // to allow users to join after the start message.
      DOM.startQuizButton().classList.add("hide");
      DOM.quizTitle().classList.remove("hide");
      DOM.quizQuestionContainer().classList.remove("hide");
      DOM.quizAnswerContainer().classList.remove("hide");

      break;
    case mqttTopics.ANSWERS:
      DOM.quizAnswer().innerText = message.payloadString;
      break;
    case mqttTopics.STATUS:
      switch (message.payloadString) {
        case mqttMessages.START:
          break;
        case mqttMessages.END:
          DOM.quizTitle().classList.add("hide");
          DOM.quizQuestionContainer().classList.add("hide");
          DOM.quizQuestion().innerText = "";
          DOM.quizAnswerContainer().classList.add("hide");
          DOM.quizAnswer().innerText = "";
          DOM.startQuizButton().classList.remove("hide");
          break;
        default:
          console.log("Unrecognized MQTT status message:", message.payloadString);
      }
      break;
    default:
      console.log("Unrecognized MQTT topic:", message.destinationName);
  }
}

/** Starts a new quiz session by publishing a start quiz message to the MQTT broker. */
export function startQuiz() {
  if (!mqttClient) {
    console.log("Failed to start quiz: MQTT client uninitialized.");
    return;
  }

  const message = new Paho.MQTT.Message(mqttMessages.START);
  message.destinationName = mqttTopics.STATUS;
  mqttClient.send(message);
}

/** Disconnects from the MQTT broker. */
export function disconnectMQTT() {
  if (mqttClient?.isConnected()) {
    mqttClient.disconnect();
    console.log("Disconnected from MQTT broker.");
  }
}
