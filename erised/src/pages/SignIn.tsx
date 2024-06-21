import { ForkOutlined } from "@ant-design/icons";
import SignInForm from "../forms/SignInForm";
import { headerLogo } from "../Styles";

function SignIn({ onSignIn }: { onSignIn: () => void }) {
  return (
    <>
      <ForkOutlined style={headerLogo} />
      <SignInForm onSignIn={onSignIn} />
    </>
  );
}

export default SignIn;
