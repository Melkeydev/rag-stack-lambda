/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState, useEffect } from "react";
import "./App.css";

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
  const [count, setCount] = useState(0);

  useEffect(() => {
    readConfig().then(() => {
      console.log("this config.apiUrl", config.apiUrl);
    });
  }, []);

  async function callTest() {
    try {
      const response = await fetch(`https://${config.apiUrl}/test`, {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`Error! status: ${response.status}`);
      }

      const result = await response.json();

      console.log("result is: ", JSON.stringify(result, null, 4));
    } catch (err: any) {
      if (err.message) {
        console.log(err.message);
      }
    }
  }

  return (
    <>
      <div>
        <a href="https://vitejs.dev" target="_blank"></a>
        <a href="https://react.dev" target="_blank"></a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <button onClick={() => callTest()}>Touch me</button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  );
}

export default App;
