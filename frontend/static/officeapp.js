/**
 * Code used from the Twilio example on "Build a Video Chat Application with Python, JavaScript and Twilio Programmable Video",
 * to connect and use simple twilio functions
 * https://www.twilio.com/blog/build-video-chat-application-python-javascript-twilio-programmable-video
 */
 const root = document.getElementById('root');
 const startQuiz = document.getElementById('start_quiz');
 const container = document.getElementById('container');
 const count = document.getElementById('count');
 const question_title = document.getElementById('Question_title');
const answer_title = document.getElementById('Answer_title');
 const questions = document.getElementById('questions');
 const answers = document.getElementById('answer');
 let connected = false;
 let room;
 let motion_time = new Date(2018, 11, 24, 10, 33, 30, 0);               //random date
 let motion = false;
 
 
 //Uses diffy library to detect motion on camera.
 var diffy = Diffy.create({
     resolution: { x: 10, y: 5 },
     sensitivity: 0.2,
     threshold: 25,
     debug: false,
     containerClassName: 'my-diffy-container',
     sourceDimensions: { w: 130, h: 100 },
     onFrame: function (matrix) {
         for(var i = 0; i < matrix.length; i++) {
             var array = matrix[i];
             if(array.some(x => x < 200)){
                 motion_time = new Date();                                
             }
         }
         if((Math.floor((new Date() - motion_time) / 1000)) < 10){
             motion = true;                                        //Set motion to true if motion has been detected in the last minute
             //console.log("Motion detected");                   
         }
         else{
             motion = false;                                      //Sets motion to false if motion has not been detected in the last minute
             //console.log("Motion not detected");   
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
         else if (message.destinationName == "TTM4115/t4/quiz/a"){
            answers.innerHTML = message.payloadString;
         }
         else if(message.destinationName == "TTM4115/t4/quiz/s"){     //Else if the message is from the start/stop quiz channel
             if(message.payloadString == "Quiz_started"){            
                 startQuiz.disabled = true;                          
             }
             else if(message.payloadString == "Quiz_ended"){
                 startQuiz.disabled = false;
                 questions.innerHTML = "";
                 answers.innerHTML = "";
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
 
 var officeconnect = window.setInterval(function(){     //Function is called every 5 seconds to check if office workers should connect or not.
     if (motion) {
         console.log("Motion detected");
         
         if(!connected){
             var username = "e"
             connect(username).then(() => {
                 startQuiz.disabled = false;
             }).catch(() => {
                 alert('Connection failed. Is the backend running?');
             });
         }
     }
     else if(!motion){
         console.log("Motion not detected");
         if(connected){
             disconnect();
             connected = false;
             startQuiz.innerHTML = 'Start Quiz';
             questions.innerHTML = '';
             answers.innerHTML = '';
             startQuiz.disabled = true;
         }
     }
 }, 5000);
 
 
 //Function taken from the twilio example to connect to the stream by using twilio
 function connect(username) {
     if(!connected){
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
     }  
 };
 
 //Function taken from the twilio example to update the participant count
 function updateParticipantCount() {
     if (!connected){
        count.innerHTML = 'You are disconnected';
        answer_title.innerText = "";
        question_title.innerHTML = "";
     }
     else{
         count.innerHTML = (room.participants.size + 1) + ' participants online.';
        answer_title.innerText = "Answer:";
        question_title.innerHTML = "Question:";
     }
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
 function startQuizHandler(event) {
     event.preventDefault();
     var message = new Paho.MQTT.Message("Quiz_started");
     message.destinationName = "TTM4115/t4/quiz/s";
     client.send(message);
     startQuiz.disabled = true;
 };

 
 addLocalVideo();
 startQuiz.addEventListener('click', startQuizHandler);

 