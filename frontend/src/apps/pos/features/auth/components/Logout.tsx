import React, { useEffect } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../hooks/useAuth";
import BaseCard from "../../../components/BaseCard";
import { Spinner } from "flowbite-react";
import { logout } from "../../../../../core/api/auth";
import { useConfig } from "@core/config/hooks/useConfig";

const Logout: React.FC = () => {
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
    <BaseCard title="">
      Logging you out...
      <br />
      <Spinner className="mt-3" />
    </BaseCard>
  );
};

export default Logout;
