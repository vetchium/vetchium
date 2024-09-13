import React from 'react';
import { Layout, Menu } from 'antd';
import { UserOutlined, FileSearchOutlined } from '@ant-design/icons';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const { Sider } = Layout;

const Sidebar: React.FC = () => {
  const { t } = useTranslation();

  return (
    <Sider width={200} className="site-layout-background">
      <Menu
        mode="inline"
        defaultSelectedKeys={['profile']}
        style={{ height: '100%', borderRight: 0 }}
      >
        <Menu.Item key="profile" icon={<UserOutlined />}>
          <Link to="/profile">{t('common.profile')}</Link>
        </Menu.Item>
        <Menu.Item key="jobs" icon={<FileSearchOutlined />}>
          <Link to="/jobs">{t('common.jobs')}</Link>
        </Menu.Item>
      </Menu>
    </Sider>
  );
};

export default Sidebar;
