import { Outlet } from "react-router-dom";
import { faBars } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

function Layout() {
  return (
    <div className={"app"}>
      {/*Navbar*/}
      <nav className="bg-[#141414] w-full h-[60px] flex flex-row items-center px-10 justify-around">
        <h1
          className={
            "hidden sm:block text-white sm:text-[22px] text-[20px] cursor-pointer"
          }
          onClick={() => {
            window.location.href = "/";
          }}
        >
          Local Cloud
        </h1>
        <input
          type={"text"}
          placeholder={"Search"}
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
          <FontAwesomeIcon icon={faBars} className={"text-2xl"} />
        </div>
        {/*For lg > md*/}
        <div className={"hidden md:block"}>
          <h1
            className={
              "px-5 py-2 bg-zinc-900 rounded-md border border-zinc-700 cursor-pointer hover:bg-zinc-700 duration-100"
            }
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
      {/*<footer className={"bg-zinc-950 text-white p-10 flex justify-center"}>Footer</footer>*/}
    </div>
  );
}

export default Layout;
