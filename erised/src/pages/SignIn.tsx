import { ForkOutlined } from "@ant-design/icons";
import SignInForm from "../forms/SignInForm";
import { headerLogo } from "../Styles";

function SignIn() {
  return (
    <>
      <ForkOutlined style={headerLogo} />
      <SignInForm />
    </>
  );
}

export default SignIn;
