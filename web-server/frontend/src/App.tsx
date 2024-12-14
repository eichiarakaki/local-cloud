import "./App.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import VideoPage from "./pages/SingleVideoPage.tsx";
import Layout from "./Layout.tsx";
import NotFound from "./pages/NotFound.tsx";
import Downloader from "./pages/Downloader.tsx";

function App() {
  return (
    <Router>
      <Routes>
        {/* Home Page */}
        <Route path={"/"} element={<Layout />}>
          <Route index element={<Home />}></Route>
          {/* Single Video Page */}
          <Route path={"/video/:videoData"} element={<VideoPage />}></Route>
        </Route>

        {/*Downloader Page*/}
        <Route path={"/download"} element={<Downloader />} />
        {/*Managing 404 requests*/}
        <Route path={""} element={<NotFound />} />
      </Routes>
    </Router>
  );
}

export default App;
