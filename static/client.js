const joinButton = document.getElementById('join');
const messagesDiv = document.getElementById('messages');
const localVideo = document.getElementById('localVideo');
// const remoteVideo = document.getElementById('remoteVideo');

let ws;
let pc;
let channel;

const configuration = {
    iceServers: [
        { urls: 'stun:stun.l.google.com:19302' }
    ]
};

joinButton.addEventListener('click', async () => {
    channel = document.getElementById('channel').value;

    await createPeerConnection();

    if (channel === '') {
        alert('Please enter a channel name');
        return;
    }

    const wsUrl = `ws://${window.location.host}/websocket?channel=${channel}`;
    ws = new WebSocket(wsUrl);

    ws.onopen = async () => {
        console.log('WebSocket connection opened');
        // displayMessage('Connected to channel: ' + channel);

        // Get user media (audio and video)
        const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
        localVideo.srcObject = stream;

        // Create PeerConnection and add tracks
        stream.getTracks().forEach(track => pc.addTrack(track, stream));
    };
    
    ws.onmessage = async (event) => {
        const message = JSON.parse(event.data);
        console.log('Received message:', message);
        if (!message) {
            return console.log('failed to parse msg')
        }

        switch (message.event) {
            case 'offer':
                await handleOfferMessage(message);
                break;
            case 'candidate':
                handleCandidateMessage(message);
                break;
            default:
                console.log('Unknown message type:', message.type);
        }
    };

    ws.onclose = () => {
        console.log('WebSocket connection closed');
        displayMessage('Disconnected from channel: ' + channel);
    };

    ws.onerror = (error) => {
        console.log('WebSocket error:', error);
        displayMessage('WebSocket error: ' + error);
    };
});

function displayMessage(message) {
    const messageElement = document.createElement('div');
    messageElement.textContent = message;
    messagesDiv.appendChild(messageElement);
}

async function createPeerConnection() {
    pc = new RTCPeerConnection(configuration);

    pc.ontrack = (event) => {
        console.log('Track received:', event.streams);

        if (event.track.kind === 'audio') {
            return
        }

        let el = document.createElement(event.track.kind)
        el.srcObject = event.streams[0];
        el.autoplay = true;
        el.controls = true;
        document.getElementById('remoteVideos').appendChild(el) // remoteVideos

        event.track.onmute = function (event) {
            el.play()
        }

        event.streams[0].onremovetrack = ({ track }) => {
            if (el.parentNode) {
                el.parentNode.removeChild(el)
              }
        }
    };

    pc.onicecandidate = (event) => {
        if (event.candidate) {
            console.log('ICE candidate:', event.candidate);
            ws.send(JSON.stringify({
                type: 'candidate',
                candidate: event.candidate,
                channel: channel
            }));
        }
    };
}

async function handleOfferMessage(message) {
    let offer = JSON.parse(message.data)
    if (!offer) {
        return console.log('failed to parse answer')
    }

    await pc.setRemoteDescription(offer);

    const answer = await pc.createAnswer();
    await pc.setLocalDescription(answer);

    ws.send(JSON.stringify({
        type: 'answer',
        answer: answer,
        channel: channel
    }));
}

function handleCandidateMessage(message) {
    let candidate = JSON.parse(message.data)
    if (!candidate) {
        return console.log('failed to parse candidate')
    }
    pc.addIceCandidate(candidate)
}