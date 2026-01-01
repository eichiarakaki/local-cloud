import { useLocation } from "react-router-dom";
import { useEffect } from "react";
import { API_BASE } from "../config";

function VideoPage() {
  const location = useLocation();

  const queryParams = new URLSearchParams(location.search);

  const title = queryParams.get("title");
  const createdAt = queryParams.get("created_at");
  const embeddedVideo = queryParams.get("embedded_video");

  useEffect(() => {
    document.title = title ? title : "Local Cloud";
  }, [title]);

  return (
    <div
      className={
        "w-full bg-[#101010] text-white flex flex-col items-center mt-10"
      }
    >
      <video
        src={`${API_BASE}/api/videos-storage/${embeddedVideo}`}
        controls
        preload="metadata"
        crossOrigin="anonymous"
        className={"lg:w-[65%] mb-5"}
      />
      <div className={"flex flex-col justify-between items-center text-center"}>
        <h2 className={"mb-5 sm:text-[25px] text-[20px]"}>{title}</h2>
        <p className={"text-zinc-400"}>{createdAt}</p>
      </div>
    </div>
  );
}

export default VideoPage;
