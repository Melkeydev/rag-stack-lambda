/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useState } from "react";
import "./App.css";
import { Register } from "./pages/Register";
import { HomePage } from "./pages/HomePage";
import { BrowserRouter, Route, Routes } from "react-router-dom";

function App() {
  const [apiUrl, setApiUrl] = useState("");

  const readConfig = async () => {
    let url = "";
    if (import.meta.env.MODE === "development") {
      // Use .env file during development
      url = import.meta.env.VITE_API_URL || "";
    } else {
      // Use config.json during production
      url = await fetch("./config.json")
        .then((response) => response.json())
        .then((json) => json.apiUrl);
    }
    setApiUrl(url);
  };

  useEffect(() => {
    // Read the config
    readConfig();
  }, []);

  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<HomePage apiUrl={apiUrl} />} />
          <Route path="/register" element={<Register apiUrl={apiUrl} />} />
        </Routes>
      </BrowserRouter>
    </>
  );
}

export default App;
