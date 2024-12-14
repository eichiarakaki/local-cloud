import React, { useState } from "react";

function Downloader() {
  const [inputValue, setInputValue] = useState<string>("");
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
            body: JSON.stringify({ url: inputValue }),
          },
        );

        if (!response.ok) {
          console.error("Error in POST request:", response.statusText);
        } else {
          const data = await response.json();
          console.log("Response from API:", data);
        }
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
    </div>
  );
}

export default Downloader;
