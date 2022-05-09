import { setUsername } from "./user.js";
import { DOM } from "./dom.js";
import { joinSession, leaveSession } from "./session.js";
import { socketOpen } from "./webrtc/socket.js";

/** Initializes the office client with a name for the stream. */
export function initializeOffice() {
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
}

/**
 * Watches for motion in the webcam, and automatically joins the stream if motion is detected.
 * Leaves the stream again after 2 minutes of inactivity.
 */
export function detectMotion() {
  // State for motion variables.
  let motionActive = false;
  let lastMotionTime = new Date(2018, 11, 24, 10, 33, 30, 0); // Random date in the past.

  // Uses Diffy.js library to detect motion from webcam.
  return Diffy.create({
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
          leaveSession();
          motionActive = false;
          console.log("Motion timed out, left call.");
        }
      } else {
        if (!motionActive && socketOpen) {
          joinSession();
          motionActive = true;
          console.log("Motion detected, joined call.");
        }
      }
    },
  });
}
