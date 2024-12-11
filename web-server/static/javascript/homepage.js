'use strict';

const videoContainer = document.getElementById("videos-container");

function videoHTML(filepath, title, created_at) {
    const encodedTitle = encodeURIComponent(filepath);
    const thumbnailPath = `/api/videos-storage/${encodedTitle}`
    return `<div class="video-container">
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
        console.log(data);
         // Handle the data from the response 
         data.forEach(element => { 
            videoContainer.innerHTML += videoHTML(element.filepath, element.filename, element.created_at); 
        }); 
    } catch (error) {
         // Handle errors that occurred during the fetch 
         console.error("There was a problem with the fetch operation:", error); 
        } 
    } // Call the function to fetch videos 
fetchVideos();