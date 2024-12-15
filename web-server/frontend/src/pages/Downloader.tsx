import React, { useEffect, useState } from "react";

interface ResponseData {
  server_status: string;
  queue_position: string;
  message: string;
}

function Downloader() {
  const [inputValue, setInputValue] = useState<string>("");
  const [data, setData] = useState<ResponseData | undefined>();

  // Changing the page title
  useEffect(() => {
    document.title = "Local Cloud | Downloader";
  }, []);

  const handleKeyDown = async (
    event: React.KeyboardEvent<HTMLInputElement>,
  ) => {
    if (event.key === "Enter" && inputValue.trim() !== "") {
      try {
        const response = await fetch(
          "http://localhost:3033/api/send-to-downloader-server",
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
        } else {
          const rawData = await response.json();
          const parsedData: ResponseData = JSON.parse(rawData);

          setData(parsedData);
        }

        setInputValue("");
      } catch (error) {
        console.error("Error sending POST request:", error);
      }
    }
  };

  return (
    <div className={"flex flex-col mt-[200px] items-center mx-5"}>
      <h1 className={"text-3xl mb-10 select-none text-center"}>
        Paste an URL from <span className={"text-red-600"}>YouTube</span>.
      </h1>

      <input
        type="text"
        placeholder={""}
        value={inputValue}
        onChange={(event) => setInputValue(event.target.value)}
        onKeyDown={handleKeyDown}
        className={
          "px-5 lg:py-5 py-4 max-w-[400px] min-w-[200px] w-full rounded-md border border-[#202020]" +
          " hover:border-zinc-700" +
          " cursor-pointer" +
          " outline-none" +
          " duration-100"
        }
      />
      <div
        className={
          "mt-10 border min-w-[200px] max-w-[600px] w-full min-h-[300px] max-h-[500px] rounded-md" +
          " border-[#202020] flex flex-col p-5 font-['Geist'] overflow-auto"
        }
      >
        {data ? (
          <>
            <div className={""}>
              <span>Server Status: {data.server_status}</span>
            </div>
            <div className={""}>
              <span>Queue Position: {data.queue_position}</span>
            </div>
            <div className={""}>
              <span>Server Message: {data.message}</span>
            </div>
          </>
        ) : (
          <span className={""}>Waiting for input...</span>
        )}
      </div>
    </div>
  );
}

export default Downloader;
