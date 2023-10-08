/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect } from "react";
import "./App.css";
import { Register } from "./pages/Register";
import { HomePage } from "./pages/HomePage";
import { BrowserRouter, Routes, Route } from "react-router-dom";

let config = {
  apiUrl: "",
};

async function readConfig() {
  if (import.meta.env.MODE === "development") {
    // Use .env file during development
    config.apiUrl = import.meta.env.VITE_API_URL || "";
  } else {
    // Use config.json during production
    config.apiUrl = await fetch("./config.json")
      .then((response) => response.json())
      .then((json) => json.apiUrl);
  }
}

function App() {
  useEffect(() => {
    // Read the config
    readConfig();
  }, []);

  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<HomePage apiUrl={config.apiUrl} />} />
          <Route
            path="/register"
            element={<Register apiUrl={config.apiUrl} />}
          />
        </Routes>
      </BrowserRouter>
    </>
  );
}

export default App;
