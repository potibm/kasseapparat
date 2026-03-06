import React from "react";
import { Button as FlowbiteButton } from "flowbite-react";
import PropTypes from "prop-types";

const Button = ({ size = "lg", children, ...props }) => {
  return (
    <FlowbiteButton size={size} {...props}>
      {children}
    </FlowbiteButton>
  );
};

Button.propTypes = {
  size: PropTypes.string,
  children: PropTypes.node.isRequired,
};

export default Button;
