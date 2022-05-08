import { setUsername } from "./user.js";
import { DOM } from "./dom.js";
import { joinCall, leaveCall } from "./main.js";

// Only starts motion detection if the appropriate client type env is passed.
if (env.CLIENT_TYPE === "office") {
  // Checks for location URL parameter.
  const officeLocation = new URLSearchParams(window.location.search).get("location");
  if (officeLocation) {
    const officeName = `Office ${officeLocation}`;
    setUsername(officeName);
    DOM.localVideoName().innerText = officeName;
  } else {
    window.alert(
      "Please provide office location as URL parameter.\n\nExample: /office?location=Oslo"
    );
  }

  // State for motion variables.
  let motionActive = false;
  let lastMotionTime = new Date(2018, 11, 24, 10, 33, 30, 0); // Random date in the past.

  // Uses Diffy.js library to detect motion from webcam.
  Diffy.create({
    resolution: { x: 10, y: 5 },
    sensitivity: 0.2,
    threshold: 25,
    debug: false,
    sourceDimensions: { w: 130, h: 100 },
    onFrame: (matrix) => {
      // Gets the current time.
      const now = new Date();

      // Checks for motion in the webcam.
      for (let i = 0; i < matrix.length; i++) {
        const array = matrix[i];
        if (array.some((x) => x < 200)) {
          lastMotionTime = now;
        }
      }

      const secondsSinceLastMotion = Math.floor(now.getSeconds() - lastMotionTime.getSeconds());

      // Checks for motion in the past two minutes, and joins/leaves the call accordingly.
      if (secondsSinceLastMotion >= 120) {
        if (motionActive) {
          leaveCall();
          motionActive = false;
          console.log("Motion timed out, left call.");
        }
      } else {
        if (!motionActive) {
          joinCall();
          motionActive = true;
          console.log("Motion detected, joined call.");
        }
      }
    },
  });
}
