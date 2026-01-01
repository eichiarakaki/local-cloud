import { Outlet } from "react-router-dom";
import { faDownload } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useState } from "react";
import Footer from "./components/Footer";

function Layout() {
  const [inputValue, setInputValue] = useState("");

  return (
    <div className={"app"}>
      {/*Navbar*/}
      <nav className="bg-[#141414] w-full h-[60px] flex flex-row items-center px-10 justify-around">
        <img
          src={"/public/home-logo.svg"}
          alt={"Logo"}
          className={"hidden sm:block text-white cursor-pointer w-[150px]"}
          onClick={() => {
            window.location.href = "/";
          }}
        />

        <input
          type={"text"}
          placeholder={"Search"}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              window.location.href = `/browser/?word=${inputValue}`;
            }
          }}
          value={inputValue}
          className={
            "px-5 py-2 rounded-lg border border-[#101010] hover:border-zinc-700 duration-100 outline-none" +
            " min-w-[100px] w-[350px]"
          }
        ></input>

        {/*For Mobiles*/}
        <div
          className={
            "md:hidden pl-6 cursor-pointer justify-center items-center flex"
          }
        >
          <FontAwesomeIcon
            icon={faDownload}
            className={"text-lg"}
            onClick={() => (window.location.href = "/download")}
          />
        </div>
        {/*For lg > md*/}
        <div className={"hidden md:block"}>
          <h1
            className={
              "px-5 py-2 bg-zinc-900 rounded-md border border-zinc-700 cursor-pointer hover:bg-zinc-700 duration-100"
            }
            onClick={() => (window.location.href = "/download")}
          >
            Download
          </h1>
        </div>
      </nav>

      {/*Main Content*/}
      <main>
        <Outlet />
      </main>

      {/*Footer*/}
      <Footer />
    </div>
  );
}

export default Layout;
