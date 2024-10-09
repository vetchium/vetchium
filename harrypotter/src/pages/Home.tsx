import { Flex, Layout } from "antd";
import { contentStyle, footerStyle, layoutStyle } from "../Styles";
import VetchiHeader from "../components/Header";
import Router from "../components/Router";
import Sidebar from "../components/Sidebar";

const { Footer, Content } = Layout;

function Home({ onSignOut }: { onSignOut: () => void }) {
  return (
    <Flex gap="middle" wrap>
      <Layout style={layoutStyle}>
        <VetchiHeader />
        <Layout>
          <Sidebar onSignOut={onSignOut} />
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
