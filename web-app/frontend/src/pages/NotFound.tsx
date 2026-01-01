import { useEffect } from "react";

function NotFound() {
  useEffect(() => {
    document.title = "Page Not Found";
  }, []);

  return <div>Not found</div>;
}

export default NotFound;
