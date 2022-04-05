/**
 * Code used from the Twilio example on "Build a Video Chat Application with Python, JavaScript and Twilio Programmable Video",
 * to connect and use simple twilio functions
 * https://www.twilio.com/blog/build-video-chat-application-python-javascript-twilio-programmable-video
 */
 const root = document.getElementById('root');
 const usernameInput = document.getElementById('username');
 const button = document.getElementById('join_leave');
 const startQuiz = document.getElementById('start_quiz');
 const container = document.getElementById('container');
 const count = document.getElementById('count');
 const questions = document.getElementById('questions');
 let connected = false;
 let room;
 
 //Uses diffy library to detect motion on camera.
 var diffy = Diffy.create({
     resolution: { x: 10, y: 5 },
     sensitivity: 0.2,
     threshold: 25,
     debug: false,
     containerClassName: 'my-diffy-container',
     sourceDimensions: { w: 130, h: 100 },
     onFrame: function (matrix) {
         var motion = false;
         for(var i = 0; i < matrix.length; i++) {
             var array = matrix[i];
             if(array.some(x => x < 200)){
                 motion = true;                                //Sets motion to true if the matrix contains values under 200
             }
         }
         if(motion){
             console.log("Motion detected")                   //TODO: stream automatically from office if motion is true
         }
      }
   });
 
 var client = new Paho.MQTT.Client('test.mosquitto.org', 8080, "5dc8b26287094cf6sa9a30kcd5eae3cd539");
 client.connect({
     onSuccess:function(){
         console.log("connected");
         client.subscribe('TTM4115/t4/quiz/#');
     }
 });
 
 //Function on what to do when a new message arrives from mqtt
 client.onMessageArrived = function(message){
     if(connected){                                                   //If the client is not connected to the stream, nothing happens
         if(message.destinationName == "TTM4115/t4/quiz/q"){          //If the message is questions
             questions.innerHTML = message.payloadString;
         }
         else if(message.destinationName == "TTM4115/t4/quiz/s"){     //Else if the message is from the start/stop quiz channel
             if(message.payloadString == "Quiz_started"){            
                 startQuiz.disabled = true;                          
             }
             else if(message.payloadString == "Quiz_ended"){
                 startQuiz.disabled = false;
                 questions.innerHTML = "-";
             }
         }
     }
     
 };
 
 
 //Function taken from the twilio example to add video
 function addLocalVideo() {
     Twilio.Video.createLocalVideoTrack().then(track => {
         let video = document.getElementById('local').firstChild;
         let trackElement = track.attach();
         trackElement.addEventListener('click', () => { zoomTrack(trackElement); });
         video.appendChild(trackElement);
     });
 };
 
 /*Function taken from the twilio example to handle the connect button.
 *It eithers connects the client to the stream if the client is not connected,
 *or it disconnects the client if it is connected
 */
 function connectButtonHandler(event) {
     event.preventDefault();
     if (!connected) {
         let username = usernameInput.value;
         if (!username) {
             alert('Enter your name before connecting');
             return;
         }
         button.disabled = true;
         button.innerHTML = 'Connecting...';
         connect(username).then(() => {
             button.innerHTML = 'Leave call';
             button.disabled = false;
             startQuiz.disabled = false;
         }).catch(() => {
             alert('Connection failed. Is the backend running?');
             button.innerHTML = 'Join call';
             button.disabled = false;
         });
     }
     else {
         disconnect();
         button.innerHTML = 'Join call';
         connected = false;
         startQuiz.innerHTML = 'Start Quiz';
         questions.innerHTML = '-';
         startQuiz.disabled = true;
     }
 };
 
 //Function taken from the twilio example to connect to the stream by using twilio
 function connect(username) {
     let promise = new Promise((resolve, reject) => {
         // get a token from the back end
         let data;
         fetch('/login', {
             method: 'POST',
             body: JSON.stringify({'username': username})
         }).then(res => res.json()).then(_data => {
             // join video call
             data = _data;
             return Twilio.Video.connect(data.token);
         }).then(_room => {
             room = _room;
             room.participants.forEach(participantConnected);
             room.on('participantConnected', participantConnected);
             room.on('participantDisconnected', participantDisconnected);
             connected = true;
             updateParticipantCount();
             resolve();
         }).catch(e => {
             console.log(e);
             reject();
         });
     });
     return promise;
 };
 
 //Function taken from the twilio example to update the participant count
 function updateParticipantCount() {
     if (!connected)
         count.innerHTML = 'Disconnected.';
     else
         count.innerHTML = (room.participants.size + 1) + ' participants online.';
 };
 
 //Function taken from the twilio example to add new participants to the stream
 function participantConnected(participant) {
     let participantDiv = document.createElement('div');
     participantDiv.setAttribute('id', participant.sid);
     participantDiv.setAttribute('class', 'participant');
 
     let tracksDiv = document.createElement('div');
     participantDiv.appendChild(tracksDiv);
 
     let labelDiv = document.createElement('div');
     labelDiv.setAttribute('class', 'label');
     labelDiv.innerHTML = participant.identity;
     participantDiv.appendChild(labelDiv);
 
     container.appendChild(participantDiv);
 
     participant.tracks.forEach(publication => {
         if (publication.isSubscribed)
             trackSubscribed(tracksDiv, publication.track);
     });
     participant.on('trackSubscribed', track => trackSubscribed(tracksDiv, track));
     participant.on('trackUnsubscribed', trackUnsubscribed);
 
     updateParticipantCount();
 };
 
 //Function taken from the twilio example to remove disconnected participants on the stream
 function participantDisconnected(participant) {
     document.getElementById(participant.sid).remove();
     updateParticipantCount();
 };
 
 //Function taken from the twilio example 
 function trackSubscribed(div, track) {
     let trackElement = track.attach();
     trackElement.addEventListener('click', () => { zoomTrack(trackElement); });
     div.appendChild(trackElement);
 };
 
 //Function taken from the twilio example 
 function trackUnsubscribed(track) {
     track.detach().forEach(element => {
         if (element.classList.contains('participantZoomed')) {
             zoomTrack(element);
         }
         element.remove()
     });
 };
 
 //Function taken from the twilio example to disconnect the client from the server and remove every other participant from the screen
 function disconnect() {
     room.disconnect();
     while (container.lastChild.id != 'local')
         container.removeChild(container.lastChild);
     button.innerHTML = 'Join call';
     startQuiz.disabled = true;
     connected = false;
     updateParticipantCount();
 };
 
 //Function taken from the twilio example to handle zoom in and zoom outs of specified participants
 function zoomTrack(trackElement) {
     if (!trackElement.classList.contains('trackZoomed')) {
         // zoom in
         container.childNodes.forEach(participant => {
             if (participant.classList && participant.classList.contains('participant')) {
                 let zoomed = false;
                 participant.childNodes[0].childNodes.forEach(track => {
                     if (track === trackElement) {
                         track.classList.add('trackZoomed')
                         zoomed = true;
                     }
                 });
                 if (zoomed) {
                     participant.classList.add('participantZoomed');
                 }
                 else {
                     participant.classList.add('participantHidden');
                 }
             }
         });
     }
     else {
         // zoom out
         container.childNodes.forEach(participant => {
             if (participant.classList && participant.classList.contains('participant')) {
                 participant.childNodes[0].childNodes.forEach(track => {
                     if (track === trackElement) {
                         track.classList.remove('trackZoomed');
                     }
                 });
                 participant.classList.remove('participantZoomed')
                 participant.classList.remove('participantHidden')
             }
         });
     }
 };
 
 //Function added to handle the start quiz button, uses mqtt broker to send the quiz started message to the other clients
 function startQuizHandler() {
     event.preventDefault();
     message = new Paho.MQTT.Message("Quiz_started");
     message.destinationName = "TTM4115/t4/quiz/s";
     client.send(message);
     startQuiz.disabled = true;
 };
 
 
 addLocalVideo();
 button.addEventListener('click', connectButtonHandler);
 startQuiz.addEventListener('click', startQuizHandler);
 
