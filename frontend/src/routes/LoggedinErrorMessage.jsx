import React from "react";
import MinimalBaseCard from "../components/MinimalBaseCard";
import { Button } from "flowbite-react";
import { Link } from "react-router-dom";

export const LoggedinErrorMessage = () => {
  return (
    <MinimalBaseCard title="Error">
      <p>This page is available for logged out users, only.</p>

      <Button as={Link} to="/logout" className="mt-4">
        Logout
      </Button>
    </MinimalBaseCard>
  );
};
