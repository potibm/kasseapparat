import React from "react";
import { SaveButton, Toolbar, ToolbarProps } from "react-admin";

export const UserEditToolbar: React.FC<ToolbarProps> = (props: any) => (
  <Toolbar {...props}>
    <SaveButton />
  </Toolbar>
);

export default UserEditToolbar;
