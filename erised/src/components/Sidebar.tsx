import { LogoutOutlined, SettingFilled } from "@ant-design/icons";
import { Layout, Menu, MenuProps } from "antd";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { siderStyle } from "../Styles";
import t from "../i18n/i18n";

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
      label: t("sidebar.openings"),
    },
    {
      key: "/candidates",
      label: t("sidebar.candidates"),
    },
    {
      key: "org-settings",
      label: t("sidebar.org_settings"),
      children: [
        {
          key: "/org-settings/users",
          label: t("sidebar.users"),
        },
        {
          key: "/org-settings/locations",
          label: t("create_opening.locations"),
        },
        {
          key: "/org-settings/departments",
          label: t("sidebar.departments"),
        },
      ],
    },
    {
      key: "/account-settings",
      label: t("sidebar.account_settings"),
      icon: <SettingFilled />,
    },
    {
      key: "/signout",
      label: t("sidebar.sign_out"),
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
        items={items}
        style={{ height: "100%" }}
      />
    </Sider>
  );
}

export default Sidebar;
