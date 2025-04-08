"use client";

import { config } from "@/config";
import MenuIcon from "@mui/icons-material/Menu";
import AppBar from "@mui/material/AppBar";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Toolbar from "@mui/material/Toolbar";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

interface HeaderProps {
  onMenuClick: () => void;
}

export default function Header({ onMenuClick }: HeaderProps) {
  const router = useRouter();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [userHandle, setUserHandle] = useState<string | null>(null);
  const [profilePicUrl, setProfilePicUrl] = useState<string | undefined>(
    undefined
  );

  useEffect(() => {
    const fetchHandle = async () => {
      try {
        const sessionToken = Cookies.get("session_token");
        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/get-my-handle`,
          {
            headers: {
              Authorization: `Bearer ${sessionToken}`,
            },
          }
        );
        if (response.ok) {
          const data = await response.json();
          setUserHandle(data.handle);
        }
      } catch (error) {
        console.error("Failed to fetch user handle:", error);
      }
    };
    fetchHandle();
  }, []);

  useEffect(() => {
    const fetchProfilePicture = async () => {
      if (!userHandle) return;

      try {
        const sessionToken = Cookies.get("session_token");
        const response = await fetch(
          `${config.API_SERVER_PREFIX}/hub/profile-picture/${userHandle}`,
          {
            headers: {
              Authorization: `Bearer ${sessionToken}`,
            },
          }
        );

        if (response.ok) {
          const blob = await response.blob();
          const url = URL.createObjectURL(blob);
          setProfilePicUrl(url);
        }
      } catch (error) {
        console.error("Failed to fetch profile picture:", error);
      }
    };

    fetchProfilePicture();

    // Cleanup function to revoke the blob URL
    return () => {
      if (profilePicUrl) {
        URL.revokeObjectURL(profilePicUrl);
      }
    };
  }, [userHandle]);

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    // Clear all authentication tokens
    Cookies.remove("session_token", { path: "/" });
    Cookies.remove("tfa_token", { path: "/" });

    // Close the menu
    handleClose();

    // Redirect to login page
    router.push("/login");
  };

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
          {/* Logo */}
          <img
            src="/logo.webp"
            alt="Vetchium Logo"
            width={60}
            height={60}
            style={{ display: "block" }}
          />
        </Box>
        <div>
          <IconButton
            size="large"
            aria-label="account of current user"
            aria-controls="menu-appbar"
            aria-haspopup="true"
            onClick={handleMenu}
            color="inherit"
            sx={{ padding: 0.5 }}
          >
            <Avatar src={profilePicUrl} sx={{ width: 40, height: 40 }} />
          </IconButton>
          <Menu
            id="menu-appbar"
            anchorEl={anchorEl}
            anchorOrigin={{
              vertical: "bottom",
              horizontal: "right",
            }}
            keepMounted
            transformOrigin={{
              vertical: "top",
              horizontal: "right",
            }}
            open={Boolean(anchorEl)}
            onClose={handleClose}
          >
            <MenuItem onClick={handleLogout}>Logout</MenuItem>
          </Menu>
        </div>
      </Toolbar>
    </AppBar>
  );
}
