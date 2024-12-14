import {useEffect, useState} from "react";
import VideoCard from "../elements/videoElement.tsx";

interface Video {
    filepath: string;
    filename: string;
    thumbnail: string;
    created_at: string;
}

function Home() {
    const [videos, setVideos] = useState<Video[]>([]);  // Almacenamos los videos

    useEffect(() => {
        fetch("/api/videos")
            .then((response) => {
                    if (!response.ok) {
                        throw new Error("Response was not ok")
                    }
                    return response.json()
                }
            ).then((json: Video[]) => setVideos(json))
            .catch((error) => console.error("Error fetching API:", error))
    }, []);

    return (
            <div className='h-screen bg-[#101010] pt-5'>
                <div className="p-10 flex flex-row flex-wrap justify-center">
                  {videos.length === 0 ? (
                    <div>Loading...</div>
                  ) : (
                    videos.map((video, index) => (
                      <VideoCard
                        key={index}
                        filepath={video.filepath}
                        filename={video.filename}
                        thumbnail={video.thumbnail}
                        created_at={video.created_at}
                      />
                    ))
                    )
                  }
                </div>
            </div>
    )
}

export default Home;