import { ForkOutlined, GlobalOutlined } from "@ant-design/icons";
import { theme, Layout } from "antd";
import { headerLogo, headerMenuIcon, headerStyle } from "../Styles";

const { Header } = Layout;
const { useToken } = theme;

function VetchiHeader() {
  const { token } = useToken();

  return (
    <Header
      style={{
        ...headerStyle,
        backgroundColor: token.colorPrimary,
      }}
    >
      <ForkOutlined style={headerLogo} />
      <GlobalOutlined style={headerMenuIcon} />
    </Header>
  );
}

export default VetchiHeader;
