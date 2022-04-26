const peers = [];

export function createPeerConnection() {
  const peer = new RTCPeerConnection({
    iceServers: [
      {
        // Open source STUN service (https://www.stunprotocol.org/)
        urls: "stun:stun.stunprotocol.org",
      },
    ],
  });

  peer.onicecandidate = handleICECandidate;
  peer.ontrack = handleTrack;
  peer.onnegotiationneeded = handleNegotiationNeeded;
  peer.onremovetrack = handleRemoveTrack;
  peer.oniceconnectionstatechange = handleICEConnectionStateChange;
  peer.onicegatheringstatechange = handleICEGatheringStateChange;
  peer.onsignalingstatechange = handleSignalingStateChange;

  peers.push(peer);
}

function handleICECandidate() {}

function handleTrack() {}

function handleNegotiationNeeded() {}

function handleRemoveTrack() {}

function handleICEConnectionStateChange() {}

function handleICEGatheringStateChange() {}

function handleSignalingStateChange() {}
