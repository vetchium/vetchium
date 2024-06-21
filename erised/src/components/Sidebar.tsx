import { LogoutOutlined, SettingFilled } from "@ant-design/icons";
import { Layout, Menu } from "antd";
import { useNavigate } from "react-router-dom";
import t from "../i18n/i18n";
import { siderStyle } from "../Styles";

const { Sider } = Layout;

function Sidebar({ onSignOut }: { onSignOut: () => void }) {
  const navigate = useNavigate();

  return (
    <Sider width="15%" style={siderStyle}>
      <Menu
        onClick={(item) => {
          if (item.key === "/signout") {
            onSignOut();
          } else {
            navigate(item.key);
          }
        }}
      >
        <Menu.Item key="/openings">{t("openings")}</Menu.Item>
        <Menu.Item key="/org-settings">Org Settings</Menu.Item>
        <Menu.Item key="/account-settings" icon={<SettingFilled />}>
          Account Settings
        </Menu.Item>
        <Menu.Item key="/signout" icon={<LogoutOutlined />}>
          Sign Out
        </Menu.Item>
      </Menu>
    </Sider>
  );
}

export default Sidebar;
