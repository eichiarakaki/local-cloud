'use strict';

const videoContainer = document.getElementById("videos-container");

function videoHTML(filepath, title, thumbnail, created_at) {
    const encodedTitle = encodeURIComponent(title);
    const encodedThumbnail = encodeURIComponent(thumbnail);
    const thumbnailPath = `/api/videos-storage/${encodedThumbnail}`

    // Creating video container
    const videoContainer = document.createElement("div");
    videoContainer.classList.add("video-container");
    videoContainer.onclick = () => (openVideo(encodedTitle, title, filepath, created_at));
    
    // Creating thumbnail
    const thumbnailElement = document.createElement("img");
    thumbnailElement.classList.add("video-thumbnail");
    thumbnailElement.src = thumbnailPath;

    // Video texts container
    const videoTextsContainer = document.createElement("div");
    videoTextsContainer.classList.add("video-texts-container");

    // Creating spans
    const titleSpan = document.createElement("span");
    const createdTimeSpan = document.createElement("span");
    titleSpan.classList.add("video-title");
    createdTimeSpan.classList.add("video-created_at");
    titleSpan.innerHTML = title;
    createdTimeSpan.innerHTML = created_at;

    // Assembling
    videoTextsContainer.appendChild(titleSpan);
    videoTextsContainer.appendChild(createdTimeSpan);
    videoContainer.append(thumbnailElement);
    videoContainer.append(videoTextsContainer);

    return videoContainer;
}

const APIPrefix = "/api";
const videosURL = `${APIPrefix}/videos`;
async function fetchVideos() { 
    try { 
        const response = await fetch(videosURL); 
        // Check if the response is ok (status code 200-299) 
        if (!response.ok) { 
            throw new Error("Network response was not ok " + response.statusText); 
        } 
        // Parse the JSON from the response 
        const data = await response.json();
         // Handle the data from the response 
         data.forEach(element => { 
            const videoElement = videoHTML(element.filepath, element.filename, element.thumbnail, element.created_at);
            videoContainer.appendChild(videoElement);
        }); 
    } catch (error) {
         // Handle errors that occurred during the fetch 
         console.error("There was a problem with the fetch operation:", error); 
        } 
    } // Call the function to fetch videos 
fetchVideos();

function openVideo(videoName, videoTitle, embedVideo, createdAt) {
    const encodedTitle = encodeURIComponent(videoTitle);
    const encodedCreatedAt = encodeURIComponent(createdAt);
    const encodedEmbeddedVideo = encodeURIComponent(embedVideo);
    const videoURL = `/video/${videoName}?title=${encodedTitle}&created_at=${encodedCreatedAt}&embedded_video=${encodedEmbeddedVideo}`;
    window.location.href = videoURL;
}