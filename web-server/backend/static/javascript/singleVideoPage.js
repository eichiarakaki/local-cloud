'use strict';

const videoContainer = document.getElementById("video-container");



function videoHTML(title, embeddedVideo, created_at) {
    // Create the video container
    const videoContainer = document.createElement("div");
    videoContainer.classList.add("video-container");

    // Creates the video element
    const videoElement = document.createElement("video");
    videoElement.src = `/api/videos-storage/${embeddedVideo}`;
    videoElement.classList.add("embed-video");
    videoElement.controls = true;

    // Creates the text container
    const textsContainer = document.createElement("div");
    textsContainer.classList.add("video-texts-container");

    // Creates the text elements
    const titleElement = document.createElement("span");
    titleElement.classList.add("video-title");
    titleElement.textContent = title;

    const createdAtElement = document.createElement("span");
    createdAtElement.classList.add("video-created_at");
    createdAtElement.textContent = created_at;

    // Assembling the elements
    textsContainer.appendChild(titleElement);
    textsContainer.appendChild(createdAtElement);

    videoContainer.appendChild(videoElement);
    videoContainer.appendChild(textsContainer);

    return videoContainer;
}

// Getting the Query parameters
const urlParams = new URLSearchParams(window.location.search);
const videoTitle = urlParams.get('title') || 'Unknown Title';
const createdAt = urlParams.get('created_at') || 'Unknown Created Date';
const embeddedVideo = urlParams.get('embedded_video');

const container = document.getElementById("video-container");
const videoElement = videoHTML(videoTitle, embeddedVideo, createdAt);
container.appendChild(videoElement);