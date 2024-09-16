import { ForkOutlined } from "@ant-design/icons";
import LogInForm from "../forms/LogInForm";
import { headerLogo } from "../Styles";

function LogIn({ onLogIn }: { onLogIn: () => void }) {
  return (
    <>
      <ForkOutlined style={headerLogo} />
      <LogInForm onLogIn={onLogIn} />
    </>
  );
}

export default LogIn;
