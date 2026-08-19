import { defaultTheme } from "react-admin";
import type { ThemeOptions } from "@mui/material";

// --- LIGHT THEME ---
export const MyTheme: ThemeOptions = {
  ...defaultTheme,
  palette: {
    mode: "light",
    primary: {
      main: "#4F46E5",
    },
    secondary: {
      main: "#ff3873",
      contrastText: "#ffffff",
    },
    error: {
      main: "#ef4444",
    },
    background: {
      default: "#f8f9fa",
      paper: "#ffffff",
    },
  },
  components: {
    ...defaultTheme.components,
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: "8px",
          fontWeight: 600,
          textTransform: "none",
        },
      },
    },
  },
};

// --- DARK THEME ---
export const MyDarkTheme: ThemeOptions = {
  ...defaultTheme,
  palette: {
    mode: "dark",
    primary: {
      main: "#6366F1",
    },
    secondary: {
      main: "#121212",
      contrastText: "#ff3873",
    },
    error: {
      main: "#ef4444",
    },
    background: {
      default: "#121212",
      paper: "#1e1e1e",
    },
  },
  components: {
    ...defaultTheme.components,
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: "8px",
          fontWeight: 600,
          textTransform: "none",
        },
      },
    },
  },
};
