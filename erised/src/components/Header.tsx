import { ForkOutlined, GlobalOutlined } from "@ant-design/icons";
import { Layout } from "antd";
import { headerLogo, headerStyle } from "../Styles";

const { Header } = Layout;

function VetchiHeader() {
  return (
    <Header style={headerStyle}>
      <ForkOutlined style={headerLogo} />
      <GlobalOutlined style={headerLogo} />
    </Header>
  );
}

export default VetchiHeader;
