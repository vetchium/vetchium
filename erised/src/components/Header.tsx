import { ForkOutlined, GlobalOutlined } from "@ant-design/icons";
import { Flex, Layout } from "antd";
import { headerLogo, headerMenuIcon, headerStyle } from "../Styles";

const { Header } = Layout;

function VetchiHeader() {
  return (
    <Header style={headerStyle}>
      <Flex justify="space-between">
        <Flex align="flex-start" justify="center">
          <ForkOutlined style={headerLogo} />
        </Flex>
        <Flex align="flex-end" justify="center">
          <GlobalOutlined style={headerMenuIcon} />
        </Flex>
      </Flex>
    </Header>
  );
}

export default VetchiHeader;
