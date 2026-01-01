"use strict";

export function VideosNotFound(title, subtitle) {
    // Creating video container
    const container = document.createElement("div");
    container.classList.add("not-found");

    // Creating spans
    const titleSpan = document.createElement("span");
    const subtitleSpan = document.createElement("span");
    titleSpan.classList.add("not-found-title");
    subtitleSpan.classList.add("not-found-subtitle")
    
    titleSpan.innerHTML = title;
    subtitleSpan.innerHTML = subtitle;

    // Assembling
    container.appendChild(titleSpan);
    container.appendChild(subtitleSpan);

    return container;
}