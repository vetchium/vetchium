import { LogoutOutlined, SettingFilled } from "@ant-design/icons";
import { Flex, Layout, Menu } from "antd";
import { useNavigate } from "react-router-dom";
import {
  contentStyle,
  footerStyle,
  headerStyle,
  layoutStyle,
  siderStyle,
} from "./Styles";
import Router from "./components/Router";

const { Header, Footer, Sider, Content } = Layout;

function App() {
  const navigate = useNavigate();

  return (
    <Flex gap="middle" wrap>
      <Layout style={layoutStyle}>
        <Header style={headerStyle}>Header</Header>
        <Layout>
          <Sider width="15%" style={siderStyle}>
            <Menu
              onClick={(item) => {
                navigate(item.key);
              }}
            >
              <Menu.Item key="/openings">Openings</Menu.Item>
              <Menu.Item key="/org-settings">Org Settings</Menu.Item>
              <Menu.Item key="/account-settings" icon={<SettingFilled />}>
                Account Settings
              </Menu.Item>
              <Menu.Item key="/signout" icon={<LogoutOutlined />}>
                Sign Out
              </Menu.Item>
            </Menu>
          </Sider>
          <Content style={contentStyle}>
            <Router />
          </Content>
        </Layout>
        <Footer style={footerStyle}>Footer</Footer>
      </Layout>
    </Flex>
  );
}

export default App;
