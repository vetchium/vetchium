import React from 'react';
import { Layout } from 'antd';
import { useTranslation } from 'react-i18next';

const { Footer: AntFooter } = Layout;

const Footer: React.FC = () => {
  const { t } = useTranslation();

  return (
    <AntFooter style={{ textAlign: 'center' }}>
      {t('common.appName')} ©{new Date().getFullYear()} {t('common.footerText')}
    </AntFooter>
  );
};

export default Footer;
