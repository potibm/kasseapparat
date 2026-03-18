import React from "react";
import Button from "../../../../components/Button";
import { ButtonProps } from "flowbite-react";

const SidebarKeyboardKey: React.FC<ButtonProps> = ({
  children,
  onClick,
  ...props
}) => {
  return (
    <Button
      onClick={onClick}
      size="md"
      className="min-w-12 max-w-12"
      {...props}
    >
      {children}
    </Button>
  );
};

export default SidebarKeyboardKey;
