import React, {ReactNode} from "react";
import { Button as FlowbiteButton, ButtonProps as FlowbiteButtonProps } from "flowbite-react";

export const Button: React.FC<FlowbiteButtonProps> = ({ size = "lg", children, ...props }) => {
  return (
    <FlowbiteButton size={size} {...props}>
      {children}
    </FlowbiteButton>
  );
};

export default Button;
