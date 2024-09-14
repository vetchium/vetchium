import { Button, Typography } from 'antd';
import { useTranslation } from 'react-i18next';

const { Title, Paragraph } = Typography;

function Home({ onSignOut }: { onSignOut: () => void }) {
  const { t } = useTranslation();

  return (
    <div>
      <Title>{t('home.welcome')}</Title>
      <Paragraph>{t('home.description')}</Paragraph>
      <Button onClick={onSignOut}>{t('home.signOut')}</Button>
    </div>
  );
}

export default Home;
