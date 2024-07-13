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
      label: t("org_settings"),
      children: [
        {
          key: "/org-settings/users",
          label: t("users"),
        },
        {
          key: "/org-settings/locations",
          label: t("locations"),
        },
        {
          key: "/org-settings/departments",
          label: t("departments"),
        },
      ],
    },
    {
      key: "/account-settings",
      label: t("account_settings"),
      icon: <SettingFilled />,
    },
    {
      key: "/signout",
      label: t("sign_out"),
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
