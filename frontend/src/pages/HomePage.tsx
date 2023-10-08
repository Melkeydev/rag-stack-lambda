import { useState } from "react";

interface HomePageProps {
  apiUrl: string;
}

export const HomePage = ({ apiUrl }: HomePageProps) => {
  const [count, setCount] = useState(0);

  async function callTest() {
    try {
      const response = await fetch(`${apiUrl}test`, {
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
      <h1>Vite + React + Tailwind</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <button className="bg-red-500" onClick={() => callTest()}>
          Test AWS
        </button>
      </div>
    </>
  );
};
