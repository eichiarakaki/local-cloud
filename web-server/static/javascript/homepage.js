'use strict';

const videoContainer = document.getElementById("videos-container");

function videoHTML(filepath, title, thumbnail, created_at) {
    const encodedTitle = encodeURIComponent(title);
    const encodedThumbnail = encodeURIComponent(thumbnail);
    const thumbnailPath = `/api/videos-storage/${encodedThumbnail}`

    return `<div onclick="openVideo('${encodedTitle}', '${title}', '${filepath}','${created_at}')" class="video-container">
                    <img src=${thumbnailPath} class="video-thumbnail"></img>
                    <div class="video-texts-container">
                        <span class="video-title">${title}</span>
                        <span class="video-created_at">${created_at}</span>
                    </div>    
                    </div>`;
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
            videoContainer.innerHTML += videoHTML(element.filepath, element.filename, element.thumbnial, element.created_at); // thumbnial, not thumnail (i gotta fix this shit)
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