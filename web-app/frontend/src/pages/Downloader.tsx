import React, { useEffect, useState } from "react";
import { API_BASE } from "../config";

interface ResponseData {
  server_status: string;
  queue_position: string;
  message: string;
}

function Downloader() {
  const [inputValue, setInputValue] = useState<string>("");
  const [data, setData] = useState<ResponseData | undefined>();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // Update the page title
  useEffect(() => {
    document.title = "Local Cloud | Downloader";
  }, []);

  const sendUrl = async () => {
    if (inputValue.trim() === "") return;

    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(
        `${API_BASE}/api/send-to-downloader-server`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ url: inputValue + "\n" }),
        },
      );

      if (!response.ok) {
        console.error("Error in POST request:", response);
        setError("Server returned an error");
      } else {
        try {
          const rawData = await response.json();

          let parsedData: ResponseData;

          if (typeof rawData === "string") {
            parsedData = JSON.parse(rawData);
          } else if (typeof rawData === "object" && rawData !== null) {
            parsedData = rawData as ResponseData;
          } else {
            throw new Error("Unexpected response format");
          }

          console.log("Final parsed data:", parsedData);

          const cleanedData: ResponseData = {
            server_status: parsedData.server_status,
            queue_position: parsedData.queue_position,
            message: parsedData.message,
          };

          setData(cleanedData);
        } catch (parseError) {
          console.error("Error parsing response:", parseError);
          setError("Failed to parse server response");
        }
      }
      setInputValue("");
    } catch (networkError) {
      console.error("Network error:", networkError);
      setError("Network error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter") {
      void sendUrl();
    }
  };

  return (
    <div className="h-screen bg-[#181818] flex flex-col items-center p-5">
      <h1 className="text-4xl font-bold text-white mt-10 mb-4 text-center drop-shadow-lg">
        Local Cloud Downloader
      </h1>
      <p className="text-lg text-gray-300 mb-10 text-center">
        Paste a URL from{" "}
        <span className="text-[#d72c2c] font-semibold">YouTube</span> and
        submit.
      </p>
      <div className="w-full max-w-md">
        <div className="flex flex-col gap-4">
          <input
            type="text"
            placeholder="Enter URL..."
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyDown={handleKeyDown}
            className="px-4 py-3 rounded-md bg-zinc-800 border border-gray-600 focus:outline-none focus:ring-1 focus:ring-[#b6b6b6] text-white placeholder-gray-400"
            disabled={isLoading}
          />
          <button
            onClick={sendUrl}
            disabled={isLoading || inputValue.trim() === ""}
            className="px-4 py-3 rounded-md bg-[#777777] text-black font-semibold hover:bg-[#eeeeee] transition-colors duration-150 shadow-md disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? "Processing..." : "Submit"}
          </button>
        </div>

        <div className="mt-8 p-6 bg-zinc-900 bg-opacity-70 rounded-md shadow-lg border border-gray-700 overflow-auto max-h-64">
          {isLoading ? (
            <p className="text-gray-200">Processing request...</p>
          ) : error ? (
            <div>
              <p className="text-yellow-400">{error}</p>
              <p className="text-gray-400 mt-2 text-sm">
                Check the console for more details.
              </p>
            </div>
          ) : data ? (
            <div className="text-gray-200 space-y-3">
              <div>
                <span className="font-semibold">Server Status:</span>{" "}
                <span className="text-[#dddddd]">{data.server_status}</span>
              </div>
              <div>
                <span className="font-semibold">Queue Position:</span>{" "}
                <span className="text-[#dddddd]">{data.queue_position}</span>
              </div>
              <div>
                <span className="font-semibold">Server Message:</span>{" "}
                <span className="text-[#dddddd]">{data.message}</span>
              </div>
            </div>
          ) : (
            <p className="text-gray-400">Waiting for input...</p>
          )}
        </div>
      </div>
    </div>
  );
}

export default Downloader;
