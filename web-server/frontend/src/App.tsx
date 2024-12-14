import { useEffect, useState } from 'react'
import './App.css'

function App() {
  const [number, setNumber] = useState<string>("");

  useEffect(() => {
    fetch("/api/test/2").then((response) => response.json())
    .then((data) => setNumber(data))
    .catch((error) => console.error("Error fetching API:", error))
  }, []);

  return (
    <div>
      <h1>Number is: {number}</h1>
    </div>
  )
}

export default App
