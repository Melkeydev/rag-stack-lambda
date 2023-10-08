import { useNavigate } from "react-router-dom";

interface HomePageProps {
  apiUrl: string;
}

export const HomePage = ({ apiUrl }: HomePageProps) => {
  let navigate = useNavigate();

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

  const navigateToRegister = () => {
    navigate("/register");
  };

  return (
    <>
      <h1>Vite + React + Tailwind</h1>
      <div className="card">
        <button className="bg-red-500" onClick={() => callTest()}>
          Test AWS Call
        </button>
        <button className="bg-indigo-300" onClick={navigateToRegister}>
          Test Register User
        </button>
      </div>
    </>
  );
};
