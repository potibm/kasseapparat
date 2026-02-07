import React, { useEffect } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../provider/AuthProvider";
import BaseCard from "../../components/BaseCard";
import { Spinner } from "flowbite-react";
import { logout } from "../hooks/Api";

const Logout = () => {
  const { removeSession, getToken } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    const handleLogout = async () => {
      await logout(await getToken()).finally(() => {
        removeSession();
        navigate("/", { replace: true });
      });
    };

    handleLogout();
  }, [removeSession, navigate, getToken]);

  return (
    <BaseCard>
      Logging you out...
      <br />
      <Spinner className="mt-3" />
    </BaseCard>
  );
};

export default Logout;
