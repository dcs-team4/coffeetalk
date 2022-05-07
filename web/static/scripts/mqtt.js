//@ts-check

import { getUsername } from "./user.js";
import { DOM } from "./dom.js";

const mqttChannels = {
  QUESTIONS: "TTM4115/t4/quiz/q",
  ANSWERS: "TTM4115/t4/quiz/a",
  STATUS: "TTM4115/t4/quiz/s",
};

const mqttMessages = {
  START: "Quiz_started",
  END: "Quiz_ended",
};

let mqtt_client;

/** Connects to the MQTT broker and sets up message listeners. */
export function connectMQTT() {
  const user = getUsername();
  if (!user.ok) {
    console.log("MQTT connection failed.");
    return;
  }

  mqtt_client = new Paho.MQTT.Client(
    env?.MQTT_HOST ?? "localhost",
    parseInt(env?.MQTT_PORT) ?? 1882,
    user.name
  );

  // Handler for new MQTT messages.
  mqtt_client.onMessageArrived = (message) => {
    console.log(`Message arrived: [${message.destinationName}] ${message.payloadString}`);

    switch (message.destinationName) {
      case mqttChannels.QUESTIONS:
        DOM.quizQuestion.innerText = message.payloadString;
        DOM.quizAnswer.innerText = "";
        DOM.startQuizButton.classList.add("hide");
        DOM.quizTitle.classList.remove("hide");
        DOM.quizQuestionTitle.classList.remove("hide");
        DOM.quizAnswerTitle.classList.remove("hide");
        break;
      case mqttChannels.ANSWERS:
        DOM.quizAnswer.innerText = message.payloadString;
        break;
      case mqttChannels.STATUS:
        switch (message.payloadString) {
          case mqttMessages.START:
            break;
          case mqttMessages.END:
            DOM.quizTitle.classList.add("hide");
            DOM.quizQuestionTitle.classList.add("hide");
            DOM.quizQuestion.innerText = "";
            DOM.quizAnswerTitle.classList.add("hide");
            DOM.quizAnswer.innerText = "";
            DOM.startQuizButton.classList.remove("hide");
            break;
          default:
            console.log(`Unrecognized MQTT message: ${message.payloadString}`);
        }
        break;
      default:
        console.log(`Unrecognized MQTT topic: ${message.destinationName}`);
    }
  };

  mqtt_client.onConnectionLost = ({ errorMessage }) => {
    console.log(`Connection to MQTT broker lost: ${errorMessage}`);
  };

  mqtt_client.connect({
    onSuccess: () => {
      console.log("Successfully connected to MQTT broker.");
      mqtt_client.subscribe("TTM4115/t4/quiz/#");
    },
    onFailure: ({ errorMessage }) => {
      console.log(`Failed to connect to MQTT broker: ${errorMessage}`);
    },
    useSSL: env?.ENV === "production",
  });
}

/** Starts a new quiz session by publishing a start quiz message to the MQTT broker. */
export function startQuiz() {
  const message = new Paho.MQTT.Message(mqttMessages.START);
  message.destinationName = mqttChannels.STATUS;
  mqtt_client.send(message);
}
