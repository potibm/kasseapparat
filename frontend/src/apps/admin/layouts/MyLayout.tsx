import {
  AppBar,
  AppBarProps,
  TitlePortal,
  Layout,
  LayoutProps,
} from "react-admin";
import { Typography, Box } from "@mui/material";
import { Logo } from "@core/logo/Logo";
import { Menu } from "./Menu";

export const MyAppBar = (props: AppBarProps) => (
  <AppBar {...props} color="secondary">
    <Box flex="1" display="flex" alignItems="center">
      <Logo
        style={{
          height: "32px",
          width: "32px",
          marginRight: "12px",
        }}
      />

      <Typography
        variant="h6"
        color="inherit"
        id="react-admin-title"
        sx={{
          fontWeight: 100,
          letterSpacing: ".1rem",
          textTransform: "uppercase",
          marginRight: "10em",
        }}
      >
        Kasseapparat&nbsp;
      </Typography>
      <TitlePortal />
    </Box>
  </AppBar>
);

export const MyLayout = (props: LayoutProps) => (
  <Layout {...props} appBar={MyAppBar} menu={Menu} />
);
