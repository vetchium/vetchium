"use client";

import {
  Box,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  useTheme,
  styled,
} from "@mui/material";
import { useRouter, usePathname } from "next/navigation";
import { useTranslation } from "@/hooks/useTranslation";
import HomeIcon from "@mui/icons-material/Home";
import WorkIcon from "@mui/icons-material/Work";

const drawerWidth = 240;

const DrawerHeader = styled("div")(({ theme }) => ({
  display: "flex",
  alignItems: "center",
  padding: theme.spacing(0, 1),
  ...theme.mixins.toolbar,
  justifyContent: "center",
}));

interface SidebarProps {
  open: boolean;
}

export default function Sidebar({ open }: SidebarProps) {
  const router = useRouter();
  const pathname = usePathname();
  const theme = useTheme();
  const { t } = useTranslation();

  const menuItems = [
    { text: t("common.home"), icon: <HomeIcon />, path: "/" },
    { text: t("common.openings"), icon: <WorkIcon />, path: "/openings" },
  ];

  return (
    <Drawer
      variant="permanent"
      sx={{
        width: open ? drawerWidth : theme.spacing(9),
        flexShrink: 0,
        "& .MuiDrawer-paper": {
          width: open ? drawerWidth : theme.spacing(9),
          transition: theme.transitions.create("width", {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen,
          }),
          overflowX: "hidden",
          whiteSpace: "nowrap",
          boxSizing: "border-box",
        },
      }}
    >
      <DrawerHeader>
        <Box
          sx={{
            width: 40,
            height: 40,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            visibility: open ? "visible" : "hidden",
          }}
        >
          <Box component="span" sx={{ fontSize: "1.5rem", fontWeight: "bold" }}>
            HP
          </Box>
        </Box>
      </DrawerHeader>
      <List>
        {menuItems.map((item) => (
          <ListItem key={item.text} disablePadding>
            <ListItemButton
              selected={pathname === item.path}
              onClick={() => router.push(item.path)}
              sx={{
                minHeight: 48,
                justifyContent: open ? "initial" : "center",
                px: 2.5,
              }}
            >
              <ListItemIcon
                sx={{
                  minWidth: 0,
                  mr: open ? 3 : "auto",
                  justifyContent: "center",
                }}
              >
                {item.icon}
              </ListItemIcon>
              <ListItemText
                primary={item.text}
                sx={{
                  opacity: open ? 1 : 0,
                  whiteSpace: "nowrap",
                }}
              />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Drawer>
  );
}
