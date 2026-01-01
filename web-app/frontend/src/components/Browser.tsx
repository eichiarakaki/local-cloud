import { useEffect, useState } from "react";
import VideoCard from "./videoElement.tsx";
import { useLocation } from "react-router-dom";

interface Video {
  filepath: string;
  filename: string;
  thumbnail: string;
  created_at: string;
}

function Browser() {
  const [videos, setVideos] = useState<Video[]>([]);
  const [filteredVideos, setFilteredVideos] = useState<Video[]>([]);

  const location = useLocation();
  const query = new URLSearchParams(location.search);
  const word = query.get("word")?.toLowerCase() || "";

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

  useEffect(() => {
    // Filtering the videos
    const filtered = videos.filter((video) =>
      video.filename.toLowerCase().includes(word),
    );
    setFilteredVideos(filtered);
  }, [videos, word]);

  return (
    <div className={"flex flex-row flex-wrap justify-center mt-10"}>
      {filteredVideos.length > 0 ? (
        filteredVideos.map((video, index) => (
          <VideoCard
            key={index}
            filepath={video.filepath}
            filename={video.filename}
            thumbnail={video.thumbnail}
            created_at={video.created_at}
          />
        ))
      ) : (
        <span
          onClick={() => (window.location.href = "/")}
          className={
            "hover:underline text-blue-500 cursor-pointer hover:text-blue-600 duration-100 "
          }
        >
          No Videos Were Found.
        </span>
      )}
    </div>
  );
}

export default Browser;
