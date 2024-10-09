import { ForkOutlined, GlobalOutlined } from "@ant-design/icons";
import { Layout } from "antd";
import { headerLogo, headerMenuIcon, headerStyle } from "../Styles";

const { Header } = Layout;

function VetchiHeader() {
  return (
    <Header style={headerStyle}>
      <ForkOutlined style={headerLogo} />
      <GlobalOutlined style={headerMenuIcon} />
    </Header>
  );
}

export default VetchiHeader;
