import React from 'react';
import { Typography } from 'antd';
import { useTranslation } from 'react-i18next';
import SignupForm from '../forms/SignupForm';

const { Title } = Typography;

const Signup: React.FC = () => {
  const { t } = useTranslation();

  return (
    <div>
      <Title level={2}>{t('signup.title')}</Title>
      <SignupForm />
    </div>
  );
};

export default Signup;
