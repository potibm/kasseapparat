import React, { useEffect } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../providers/auth-provider";
import BaseCard from "../../../components/BaseCard";
import { Spinner } from "flowbite-react";
import { logout } from "../hooks/api";
import { useConfig } from "../../../../../core/config/providers/config-provider";

const Logout = () => {
  const { removeSession } = useAuth();
  const navigate = useNavigate();
  const apiHost = useConfig().apiHost;

  useEffect(() => {
    const handleLogout = async () => {
      await logout(apiHost).finally(() => {
        removeSession();
        navigate("/", { replace: true });
      });
    };

    handleLogout();
  }, [removeSession, navigate, apiHost]);

  return (
    <BaseCard>
      Logging you out...
      <br />
      <Spinner className="mt-3" />
    </BaseCard>
  );
};

export default Logout;
