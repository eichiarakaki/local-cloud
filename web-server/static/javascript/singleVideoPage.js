'use strict';

const videoContainer = document.getElementById("video-container");

function videoHTML(title, embeddedVideo, created_at) {
    return `<div class="video-container">
                    <video src="/api/videos-storage/${embeddedVideo}" class="embed-video" controls></video>
                    <div class="video-texts-container">
                        <span class="video-title">${title}</span>
                        <span class="video-created_at">${created_at}</span>
                    </div>    
                    </div>`;
}

// Getting the Query parameters
const urlParams = new URLSearchParams(window.location.search);
const videoTitle = urlParams.get('title') || 'Unknown Title';
const createdAt = urlParams.get('created_at') || 'Unknown Created Date';
const embeddedVideo = urlParams.get('embedded_video');

videoContainer.innerHTML += videoHTML(videoTitle, embeddedVideo, createdAt);