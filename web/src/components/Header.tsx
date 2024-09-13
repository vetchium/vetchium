import React from 'react';
import { Layout, Menu } from 'antd';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const { Header: AntHeader } = Layout;

const Header: React.FC = () => {
  const { t } = useTranslation();

  return (
    <AntHeader>
      <div className="logo" />
      <Menu theme="dark" mode="horizontal" defaultSelectedKeys={['home']}>
        <Menu.Item key="home">
          <Link to="/">{t('common.home')}</Link>
        </Menu.Item>
        <Menu.Item key="signup">
          <Link to="/signup">{t('common.signup')}</Link>
        </Menu.Item>
        <Menu.Item key="login">
          <Link to="/login">{t('common.login')}</Link>
        </Menu.Item>
      </Menu>
    </AntHeader>
  );
};

export default Header;
