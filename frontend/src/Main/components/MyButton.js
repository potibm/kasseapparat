import React from "react";
import { Button } from "flowbite-react";
import PropTypes from "prop-types";

const MyButton = ({ size = "lg", children, ...props }) => {
  return (
    <Button size={size} {...props}>
      {children}
    </Button>
  );
};

MyButton.propTypes = {
  size: PropTypes.string,
  children: PropTypes.node.isRequired,
};

export default MyButton;
