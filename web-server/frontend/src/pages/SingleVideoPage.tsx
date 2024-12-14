import { useLocation } from "react-router-dom";

function VideoPage() {
  const location = useLocation();

  const queryParams = new URLSearchParams(location.search);

  const title = queryParams.get("title");
  const createdAt = queryParams.get("created_at");
  const embeddedVideo = queryParams.get("embedded_video");

  return (
    <div
      className={
        "w-full bg-[#101010] text-white flex flex-col items-center mt-10"
      }
    >
      <video
        src={"/api/videos-storage/" + embeddedVideo}
        controls
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
