import React from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../provider/AuthProvider";
import AuthCard from "./AuthCard";
import { Spinner } from "flowbite-react";

const Logout = () => {
  const { setToken } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    setToken();
    navigate("/", { replace: true });
  };

  setTimeout(() => {
    handleLogout();
  }, 3 * 1000);

  return (
    <AuthCard>
      Logging you out...
      <Spinner />
    </AuthCard>
  );
};

export default Logout;
