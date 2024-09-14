import { Button, Typography } from "antd";
import { useTranslation } from "react-i18next";

const { Title } = Typography;

function Login({ onSignIn }: { onSignIn: () => void }) {
  const { t } = useTranslation();

  return (
    <div>
      <Title level={2}>{t("login.title")}</Title>
      <Button onClick={onSignIn}>{t("login.signIn")}</Button>
      {/* Add login form here */}
    </div>
  );
}

export default Login;
