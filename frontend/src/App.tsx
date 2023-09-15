/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState } from "react";
// import reactLogo from "./assets/react.svg";
// import viteLogo from "/vite.svg";
import "./App.css";

let config = {
  apiUrl: "",
};

async function readConfig() {
  config = await fetch("./config.json").then((response) => response.json());
}

readConfig();

function App() {
  const [count, setCount] = useState(0);

  async function callTest() {
    console.log("hi");

    try {
      const response = await fetch(config.apiUrl, {
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
