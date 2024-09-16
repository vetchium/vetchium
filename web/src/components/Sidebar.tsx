import { LogoutOutlined, SettingFilled } from "@ant-design/icons";
import { Layout, Menu, MenuProps } from "antd";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { siderStyle } from "../Styles";
import t from "../i18n/i18n";

const { Sider } = Layout;

function Sidebar({ onLogOut }: { onLogOut: () => void }) {
  const navigate = useNavigate();

  useEffect(() => {
    if (window.location.pathname === "/") {
      navigate("/home");
    }
  }, []);

  type MenuItem = Required<MenuProps>["items"][number];
  const items: MenuItem[] = [
    {
      key: "/home",
      label: t("sidebar.home"),
    },
    {
      key: "/my-applications",
      label: t("sidebar.my_applications"),
    },
    {
      key: "/interviews",
      label: t("sidebar.interviews"),
    },
    {
      key: "/account-settings",
      label: t("sidebar.account_settings"),
      icon: <SettingFilled />,
    },
    {
      key: "/logout",
      label: t("sidebar.log_out"),
      icon: <LogoutOutlined />,
    },
  ];

  return (
    <Sider width="20%" style={siderStyle}>
      <Menu
        onClick={(item) => {
          if (item.key === "/logout") {
            onLogOut();
          } else {
            navigate(item.key);
          }
        }}
        defaultSelectedKeys={["/home"]}
        defaultOpenKeys={["org-settings"]}
        mode="inline"
        items={items}
        style={{ height: "100%" }}
      />
    </Sider>
  );
}

export default Sidebar;
