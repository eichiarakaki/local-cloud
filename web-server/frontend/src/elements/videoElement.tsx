interface VideoProps {
    filepath: string;
    filename: string;
    thumbnail: string;
    created_at: string;
  }
  
  function VideoCard({ filepath, filename, thumbnail, created_at }: VideoProps) {
    const encodedTitle = encodeURIComponent(filename);
    const encodedThumbnail = encodeURIComponent(thumbnail);
    const thumbnailPath = `/api/videos-storage/${encodedThumbnail}`;
  
    function openVideo() {
        window.location.href = `/video/${encodedTitle}?title=${encodedTitle}&created_at=${encodeURIComponent(
          created_at
      )}&embedded_video=${encodeURIComponent(filepath)}`;
    }
  
    return (
      <div className="w-[350px] m-4 cursor-pointer hover:bg-[#151515] duration-100 rounded-lg" onClick={openVideo}>
        <img className="rounded-lg" src={thumbnailPath} alt={filename} />
        <div className="flex flex-col justify-between min-h-[100px] text-center pt-5">
          <span className="text-white">{filename}</span>
          <span className="text-zinc-500 text-sm">{created_at}</span>
        </div>
      </div>
    );
  }
  
  export default VideoCard;
  