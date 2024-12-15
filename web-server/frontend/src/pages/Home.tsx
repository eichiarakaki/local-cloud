import { useEffect, useState } from "react";
import VideoCard from "../elements/videoElement.tsx";

interface Video {
  filepath: string;
  filename: string;
  thumbnail: string;
  created_at: string;
}

function Home() {
  const [videos, setVideos] = useState<Video[]>([]);

  useEffect(() => {
    // Changing the page title.
    document.title = "Home - Local Cloud";

    fetch("/api/videos")
      .then((response) => {
        if (!response.ok) {
          throw new Error("Response was not ok");
        }
        return response.json();
      })
      .then((json: Video[]) => setVideos(json))
      .catch((error) => console.error("Error fetching API:", error));
  }, []);

  return (
    <div className="bg-[#101010] pt-5">
      <div className="p-10 flex flex-row flex-wrap justify-center">
        {videos.length === 0 ? (
          <div
            className={"text-center"}
            onClick={() => (window.location.href = "/")}
          >
            <span
              className={
                "hover:underline text-blue-500 cursor-pointer hover:text-blue-600 duration-100 "
              }
            >
              No Videos Available.
            </span>
          </div>
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
        )}
      </div>
    </div>
  );
}

export default Home;
