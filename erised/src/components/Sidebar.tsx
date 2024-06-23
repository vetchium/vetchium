import { LogoutOutlined, SettingFilled } from "@ant-design/icons";
import { Layout, Menu, MenuProps } from "antd";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import t from "../i18n/i18n";
import { siderStyle } from "../Styles";

const { Sider } = Layout;

function Sidebar({ onSignOut }: { onSignOut: () => void }) {
  const navigate = useNavigate();

  useEffect(() => {
    if (window.location.pathname === "/") {
      navigate("/openings");
    }
  }, []);

  type MenuItem = Required<MenuProps>["items"][number];
  const items: MenuItem[] = [
    {
      key: "/openings",
      label: t("openings"),
    },
    {
      key: "org-settings",
      label: "Org Settings",
      children: [
        {
          key: "/org-settings/users",
          label: "Users",
        },
        {
          key: "/org-settings/locations",
          label: "Locations",
        },
      ],
    },
    {
      key: "/account-settings",
      label: "Account Settings",
      icon: <SettingFilled />,
    },
    {
      key: "/signout",
      label: "Sign Out",
      icon: <LogoutOutlined />,
    },
  ];

  return (
    <Sider width="20%" style={siderStyle}>
      <Menu
        onClick={(item) => {
          if (item.key === "/signout") {
            onSignOut();
          } else {
            navigate(item.key);
          }
        }}
        defaultSelectedKeys={["/openings"]}
        defaultOpenKeys={["org-settings"]}
        mode="inline"
        inlineCollapsed={false}
        items={items}
        style={{ height: "100%" }}
      />
    </Sider>
  );
}

export default Sidebar;
