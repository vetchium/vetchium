"use client";

import { AppBar, Toolbar, IconButton, Typography, Box } from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";

interface HeaderProps {
  onMenuClick: () => void;
}

export default function Header({ onMenuClick }: HeaderProps) {
  return (
    <AppBar
      position="fixed"
      sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}
    >
      <Toolbar>
        <IconButton
          color="inherit"
          aria-label="open drawer"
          edge="start"
          onClick={onMenuClick}
          sx={{ mr: 2 }}
        >
          <MenuIcon />
        </IconButton>
        <Box sx={{ flexGrow: 1, display: "flex", alignItems: "center" }}>
          <Box
            sx={{
              width: 40,
              height: 40,
              bgcolor: "white",
              borderRadius: 1,
              mr: 2,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <Typography sx={{ color: "primary.main", fontWeight: "bold" }}>
              HP
            </Typography>
          </Box>
          <Typography variant="h6" noWrap component="div">
            HarryPotter
          </Typography>
        </Box>
      </Toolbar>
    </AppBar>
  );
}
