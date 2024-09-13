import React from 'react';
import { Typography } from 'antd';
import { useTranslation } from 'react-i18next';

const { Title } = Typography;

const Login: React.FC = () => {
  const { t } = useTranslation();

  return (
    <div>
      <Title level={2}>{t('login.title')}</Title>
      {/* Add login form here */}
    </div>
  );
};

export default Login;
