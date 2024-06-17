import { Flex, Layout } from "antd";
import { contentStyle, footerStyle, headerStyle, layoutStyle } from "../Styles";
import Router from "../components/Router";
import Sidebar from "../components/Sidebar";

const { Header, Footer, Content } = Layout;

function Home() {
  return (
    <Flex gap="middle" wrap>
      <Layout style={layoutStyle}>
        <Header style={headerStyle}>Header</Header>
        <Layout>
          <Sidebar />
          <Content style={contentStyle}>
            <Router />
          </Content>
        </Layout>
        <Footer style={footerStyle}>Footer</Footer>
      </Layout>
    </Flex>
  );
}

export default Home;
