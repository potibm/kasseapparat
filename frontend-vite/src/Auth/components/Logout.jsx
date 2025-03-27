import React, { useEffect } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../provider/AuthProvider";
import BaseCard from "../../components/BaseCard";
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

    return () => clearTimeout(timer);
  }, [setToken, navigate]);

  return (
    <BaseCard>
      Logging you out...
      <br />
      <Spinner className="mt-3" />
    </BaseCard>
  );
};

export default Logout;
