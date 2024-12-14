function Downloader() {
  return (
    <div className={"flex flex-col mt-[200px] items-center text-center mx-5"}>
      <h1 className={"text-3xl mb-10 select-none"}>
        Paste an URL from <span className={"text-red-600"}>YouTube</span>.
      </h1>
      <input
        type="text"
        placeholder={""}
        className={
          "px-5 lg:py-6 py-4 max-w-[400px] min-w-[200px] w-full rounded-md border border-[#202020]" +
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
