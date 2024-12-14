import React, { useState } from "react";

interface ResponseData {
  message: string;
}

function Downloader() {
  const [inputValue, setInputValue] = useState<string>("");
  const [data, setData] = useState<ResponseData | undefined>();

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
          const data = await response.json();
          setData(data);
        }

        setInputValue("");
      } catch (error) {
        console.error("Error sending POST request:", error);
      }
    }
  };

  return (
    <div className={"flex flex-col mt-[200px] items-center text-center mx-5"}>
      <h1 className={"text-3xl mb-10 select-none"}>
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
      <div className={"mt-10"}>
        {data ? (
          <span className={"text-green-600"}>{data.message}</span>
        ) : (
          "Waiting for input..."
        )}
      </div>
    </div>
  );
}

export default Downloader;
