import React from "react";
import { Form, Input, Button, Card, Typography, message } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";
import styled from "styled-components";
import { useAuth } from "@/hooks/useAuth";
import { SignInCredentials } from "@/types/auth";

const { Title } = Typography;

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f0f2f5;
`;

const StyledCard = styled(Card)`
  width: 100%;
  max-width: 400px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
`;

const Logo = styled.div`
  text-align: center;
  margin-bottom: 24px;
  font-size: 24px;
  font-weight: bold;
  color: #00b96b;
`;

const SignIn: React.FC = () => {
  const { login, loading } = useAuth();

  const onFinish = async (values: SignInCredentials) => {
    try {
      await login(values);
    } catch (error) {
      message.error("Failed to sign in. Please check your credentials.");
    }
  };

  return (
    <Container>
      <StyledCard>
        <Logo>Vetchi Jobs</Logo>
        <Title level={3} style={{ textAlign: "center", marginBottom: 32 }}>
          Sign In
        </Title>
        <Form
          name="signin"
          onFinish={onFinish}
          autoComplete="off"
          layout="vertical"
        >
          <Form.Item
            name="email"
            rules={[
              { required: true, message: "Please input your email!" },
              { type: "email", message: "Please enter a valid email!" },
            ]}
          >
            <Input prefix={<UserOutlined />} placeholder="Email" size="large" />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: "Please input your password!" }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="Password"
              size="large"
            />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              size="large"
              block
              loading={loading}
            >
              Sign In
            </Button>
          </Form.Item>
        </Form>
      </StyledCard>
    </Container>
  );
};

export default SignIn;
