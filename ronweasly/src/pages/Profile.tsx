import React from "react";
import { Typography } from "antd";
import { useTranslation } from "react-i18next";

const { Title } = Typography;

const Profile: React.FC = () => {
  const { t } = useTranslation();

  return (
    <div>
      <Title level={2}>{t("profile.title")}</Title>
      {/* Add profile content here */}
    </div>
  );
};

export default Profile;
