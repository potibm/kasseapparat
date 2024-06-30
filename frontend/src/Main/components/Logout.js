import React, { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../provider/AuthProvider";
import AuthCard from "./AuthCard";
import { Spinner } from "flowbite-react";

const Logout = () => {
  const { setToken } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    const handleLogout = () => {
      setToken();
      navigate("/", { replace: true });
    };

    const timer = setTimeout(() => {
      handleLogout();
    }, 3 * 1000);

    // Bereinigungsfunktion, um den Timer zu löschen, falls die Komponente vor Ablauf des Timers demontiert wird
    return () => clearTimeout(timer);
  }, [setToken, navigate]);

  return (
    <AuthCard>
      Logging you out...
      <br />
      <Spinner className="mt-3" />
    </AuthCard>
  );
};

export default Logout;
