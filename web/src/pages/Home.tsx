import React from "react";
import { Typography } from "antd";
import { useTranslation } from "react-i18next";

const { Title, Paragraph } = Typography;

const Home: React.FC = () => {
  const { t } = useTranslation();

  return (
    <div>
      <Title>{t("home.welcome")}</Title>
      <Paragraph>{t("home.description")}</Paragraph>
    </div>
  );
};

export default Home;
