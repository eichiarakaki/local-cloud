# Local Cloud

**Local Cloud** is a microservice-based system designed to allow users within the same **Local Area Network (LAN)** to download videos via shared URLs and watch them through a locally hosted web application. 

> The system efficiently handles video downloads, stores the data, and renders them on a user-friendly front-end.


## Quick Start
> Local host is easy to configure, you can configure ports, sockets and folder in `config.json`.
- Do not forget to specify the directory where you want the downloaded videos to be stored.
- Dependencies
  - `MySQL` running on the port specified in `config.json`.
    - Make sure to create a user within the MySQL and same password as the one in `config.json`. 
    - Make sure to create a database with the same name you put in `config.json`.
    - **Do not create the TABLE**, as it is automatically created by the Downloader Server.
  - `yt-dlp` for the Downloader Server: You can install it with `pip install yt-dlp`.
  - `ffmpeg` for the Downloader Server: You can install it with your package manager.
  - `nodejs` for the Front-End: You can install it with your package manager.
- Once you have those dependencies on your system, run `go run .` inside each of these directories in order: downloader-server/ -> message-queue/ -> web-app/backend/ and finally, inside the web-server/frotend/ folder install the dependencies with npm install and then run the frontend.
## Previews
Home Page
![image](https://github.com/user-attachments/assets/fb5f8b1b-9f6c-4259-a6fd-5f55178f1573)
Download Page
![image](https://github.com/user-attachments/assets/b64a45d0-f6d0-49ff-bae3-edad9ca8e7af)


## **Key Features**

### **Web Server**  
- **Content Display**: Shows a list of videos stored in the database with the following metadata:  
  - File name  
  - File path (MP4)  
  - Automatically generated thumbnail
- **URL Submission**: Users can submit video URLs through an input field on the interface.

### **Downloader Server (Microservice)**
- **Download Management**: Handles video download requests, saves the video files, and generates related metadata.  
- **Database Updates**: After completing a download:
  - Saves the video file to the dedicated folder.
  - Generates and stores a thumbnail.
  - Updates the database with file path and metadata.  
- **TCP/IP Communication**: Maintains efficient communication with the Web Server and the Message Queue.

### **Message Queue**  
- **Inter-Service Communication**: Acts as an intermediary for asynchronous communication between the Web Server and the Downloader Server.  
- **Request Management**: Ensures that multiple download requests are processed efficiently and sequentially.

### **Database**  
- **Metadata Storage**: Stores relevant information about downloaded videos, including:  
  - File path (MP4)  
  - File name  
  - Thumbnail path  
- **Synchronization**: Ensures the front-end remains up-to-date with the downloaded videos.

### **Videos Folder**  
- **Local Storage**: All downloaded video files are saved in a designated directory accessible by the Downloader Server.


## **System Architecture**

Below is an overview of the **Local Cloud** system architecture:

![image](https://github.com/user-attachments/assets/07858769-96fc-4783-90f0-f4119d6d36d4)

### **Workflow**  
1. Users access the **Web Server** through a browser and submit a video URL.  
2. The Backend of the Web Server sends the URL to the **Message Queue** over a TCP/IP connection.  
3. The **Downloader Server** listens for messages from the queue and processes the video request.  
4. Once downloaded, the **Downloader Server** generates a thumbnail, stores the video file in the **Videos Folder**, and updates the **Database**.  
5. The **Web Server** queries the database to fetch video metadata and displays the available videos on the front-end.  


## **Technologies Used**

- **Backend**:  
  - Go (Golang) for the Web Server, Message Queue and Downloader Server.  
  - TCP/IP for inter-service communication.  
- **Frontend**:  
  - React (Vite) and Tailwind for building the user interface.  
- **Database**:  
  - MySQL (configurable).  
- **File Storage**:  
  - Local file system for a simple folder.  



