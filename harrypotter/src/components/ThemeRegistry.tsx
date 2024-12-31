"use client";

import { ThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import { createTheme } from "@mui/material";

const theme = createTheme({
  palette: {
    primary: {
      main: "#1B5E20", // Dark green
      light: "#4C8C4A",
      dark: "#003300",
      contrastText: "#ffffff",
    },
    secondary: {
      main: "#2E7D32", // Forest green
      light: "#60ad5e",
      dark: "#005005",
      contrastText: "#ffffff",
    },
    background: {
      default: "#fafafa",
      paper: "#ffffff",
    },
    text: {
      primary: "rgba(0, 0, 0, 0.87)", // Standard dark text
      secondary: "rgba(0, 0, 0, 0.6)", // Standard secondary text
    },
    error: {
      main: "#d32f2f",
    },
    warning: {
      main: "#ED6C02",
    },
    info: {
      main: "#0288d1",
    },
    success: {
      main: "#2e7d32",
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
    h1: {
      fontWeight: 500,
    },
    h2: {
      fontWeight: 500,
    },
    h3: {
      fontWeight: 500,
    },
    h4: {
      fontWeight: 500,
    },
    h5: {
      fontWeight: 500,
    },
    h6: {
      fontWeight: 500,
    },
  },
  components: {
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundColor: "#1B5E20",
        },
      },
    },
    MuiDrawer: {
      styleOverrides: {
        paper: {
          backgroundColor: "#ffffff",
          borderRight: "1px solid #e0e0e0",
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 4,
          textTransform: "none",
        },
      },
    },
    MuiListItemButton: {
      styleOverrides: {
        root: {
          "&.Mui-selected": {
            backgroundColor: "#E8F5E9",
            "&:hover": {
              backgroundColor: "#C8E6C9",
            },
          },
        },
      },
    },
  },
});

export default function ThemeRegistry({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      {children}
    </ThemeProvider>
  );
}
