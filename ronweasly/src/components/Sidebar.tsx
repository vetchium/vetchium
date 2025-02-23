"use client";

import { usePathname, useRouter } from "next/navigation";
import Link from "next/link";
import Drawer from "@mui/material/Drawer";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import HomeIcon from "@mui/icons-material/Home";
import SearchIcon from "@mui/icons-material/Search";
import AssignmentIcon from "@mui/icons-material/Assignment";
import FolderSpecialIcon from "@mui/icons-material/FolderSpecial";
import PersonIcon from "@mui/icons-material/Person";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import { styled } from "@mui/material/styles";
import { useTranslation } from "@/hooks/useTranslation";
import Cookies from "js-cookie";

const drawerWidth = 240;

const DrawerHeader = styled("div")(({ theme }) => ({
  display: "flex",
  alignItems: "center",
  padding: theme.spacing(0, 1),
  ...theme.mixins.toolbar,
  justifyContent: "flex-end",
}));

interface SidebarProps {
  open: boolean;
}

export default function Sidebar({ open }: SidebarProps) {
  const pathname = usePathname();
  const router = useRouter();
  const { t } = useTranslation();

  const handleProfileClick = (e: React.MouseEvent) => {
    e.preventDefault();
    const token = Cookies.get("session_token");
    if (!token) {
      router.push("/login");
      return;
    }
    router.push("/my-profile");
  };

  const menuItems = [
    { text: "home", icon: <HomeIcon />, path: "/" },
    { text: "findOpenings", icon: <SearchIcon />, path: "/find-openings" },
    {
      text: "myApplications",
      icon: <AssignmentIcon />,
      path: "/my-applications",
    },
    {
      text: "myCandidacies",
      icon: <FolderSpecialIcon />,
      path: "/my-candidacies",
    },
    {
      text: "myApprovals",
      icon: <CheckCircleIcon />,
      path: "/my-approvals",
    },
    {
      text: "myProfile",
      icon: <PersonIcon />,
      path: "#",
      onClick: handleProfileClick,
    },
  ];

  return (
    <Drawer
      variant="permanent"
      sx={{
        width: open ? drawerWidth : 72,
        flexShrink: 0,
        "& .MuiDrawer-paper": {
          width: open ? drawerWidth : 72,
          boxSizing: "border-box",
          transition: (theme) =>
            theme.transitions.create("width", {
              easing: theme.transitions.easing.sharp,
              duration: theme.transitions.duration.enteringScreen,
            }),
          overflowX: "hidden",
        },
      }}
    >
      <DrawerHeader />
      <List>
        {menuItems.map((item) => (
          <ListItem key={item.text} disablePadding>
            {item.onClick ? (
              <ListItemButton
                onClick={item.onClick}
                selected={pathname === "/my-profile"}
                sx={{
                  minHeight: 48,
                  justifyContent: open ? "initial" : "center",
                  px: 2.5,
                  textDecoration: "none",
                  width: "100%",
                  color: "inherit",
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
                  primary={t(`navigation.${item.text}`)}
                  sx={{
                    opacity: open ? 1 : 0,
                    whiteSpace: "nowrap",
                  }}
                />
              </ListItemButton>
            ) : (
              <Link
                href={item.path}
                style={{
                  textDecoration: "none",
                  width: "100%",
                  color: "inherit",
                }}
              >
                <ListItemButton
                  selected={pathname === item.path}
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
                    primary={t(`navigation.${item.text}`)}
                    sx={{
                      opacity: open ? 1 : 0,
                      whiteSpace: "nowrap",
                    }}
                  />
                </ListItemButton>
              </Link>
            )}
          </ListItem>
        ))}
      </List>
    </Drawer>
  );
}
