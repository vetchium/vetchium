import React, { useState } from "react";
import { Layout, Menu, Button, theme } from "antd";
import { Outlet, useNavigate, useLocation } from "react-router-dom";
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  DashboardOutlined,
  FileSearchOutlined,
  FileTextOutlined,
  UserOutlined,
  SettingOutlined,
  LogoutOutlined,
} from "@ant-design/icons";
import styled from "styled-components";
import { useAuth } from "@/hooks/useAuth";

const { Header, Sider, Content } = Layout;

const StyledLayout = styled(Layout)`
  min-height: 100vh;
`;

const Logo = styled.div`
  height: 32px;
  margin: 16px;
  color: white;
  font-size: 20px;
  font-weight: bold;
  text-align: center;
`;

const MainLayout: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const { logout } = useAuth();
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  const menuItems = [
    {
      key: "/",
      icon: <DashboardOutlined />,
      label: "Dashboard",
    },
    {
      key: "/openings",
      icon: <FileSearchOutlined />,
      label: "Find Jobs",
    },
    {
      key: "/applications",
      icon: <FileTextOutlined />,
      label: "My Applications",
    },
    {
      key: "/candidacies",
      icon: <UserOutlined />,
      label: "My Candidacies",
    },
    {
      key: "/settings",
      icon: <SettingOutlined />,
      label: "Settings",
    },
  ];

  return (
    <StyledLayout>
      <Sider trigger={null} collapsible collapsed={collapsed}>
        <Logo>Vetchi Jobs</Logo>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ padding: 0, background: colorBgContainer }}>
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{
              fontSize: "16px",
              width: 64,
              height: 64,
            }}
          />
          <Button
            type="text"
            icon={<LogoutOutlined />}
            onClick={logout}
            style={{
              fontSize: "16px",
              width: 64,
              height: 64,
              float: "right",
            }}
          />
        </Header>
        <Content
          style={{
            margin: "24px 16px",
            padding: 24,
            minHeight: 280,
            background: colorBgContainer,
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </StyledLayout>
  );
};

export default MainLayout;
