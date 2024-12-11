'use strict';

const videoContainer = document.getElementById("videos-container");

function videoHTML(title) {
    const encodedTitle = encodeURIComponent(title);
    const thumbnailPath = `/api/videos-storage/${encodedTitle}`
    return `<div class="video-container">
                    <img src=${thumbnailPath} class="thumbnail"></img>
                    <span class="title">${title}</span>
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
            videoContainer.innerHTML += videoHTML(element.filepath); 
        }); 
    } catch (error) {
         // Handle errors that occurred during the fetch 
         console.error("There was a problem with the fetch operation:", error); 
        } 
    } // Call the function to fetch videos 
fetchVideos();