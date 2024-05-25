import { useNavigate } from "react-router-dom";
import { useAuth } from "../../provider/authProvider";
import { Card, Label, Button, TextInput, Alert } from "flowbite-react";
import { getJwtToken } from "../hooks/Api";
import { useState } from "react";

const API_HOST = process.env.REACT_APP_API_HOST;

const Login = () => {
  
    const [error, setError] = useState(null);

    const { setToken, setUsername, setExpiryDate } = useAuth();
    const navigate = useNavigate();
  
    const handleLogin = (event) => {
      event.preventDefault();
      
      const username = event.target.username.value;
      const password = event.target.password.value;

      getJwtToken(API_HOST, username, password).then((auth) => {
        const token = auth.token;
        const expiryDate = auth.expire;
        console.log("Token: " + token + " Expiry: " + expiryDate + " Username: " + username);
        setToken(token);
        setUsername(username);
        setExpiryDate(expiryDate);
        navigate("/", { replace: true });
      }).catch((error) => {
        setError("Invalid username or password. Please try again.");
      })
    }

    return <div className="flex justify-center items-center h-screen">
    <Card className="max-w-sm ">
      <h5 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">Kasseapparat</h5>

      {error && <Alert color="failure">{error}</Alert>}

    <form className="flex flex-col gap-4" onSubmit={handleLogin}>
      <div>
        <div className="mb-2 block">
          <Label htmlFor="username" value="Your username" />
        </div>
        <TextInput id="username" type="text" placeholder="Username" required />
      </div>
      <div>
        <div className="mb-2 block">
          <Label htmlFor="password" value="Your password" />
        </div>
        <TextInput id="password" type="password" required />
      </div>
      <Button type="submit">Login</Button>
    </form>
  </Card></div>;
  };
  
  export default Login;