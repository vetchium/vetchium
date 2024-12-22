import React, { useState } from "react";
import { Form, Input, Button, Card, Typography, Alert } from "antd";
import { UserOutlined, LockOutlined, GlobalOutlined } from "@ant-design/icons";
import styled from "styled-components";
import { useAuth } from "@/hooks/useAuth";
import { SignInCredentials, TFARequest } from "@/types/auth";

const { Title } = Typography;

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f0f2f5;
`;

const StyledCard = styled(Card)`
  width: 100%;
  max-width: 400px;
`;

const Logo = styled.div`
  text-align: center;
  margin-bottom: 24px;
  font-size: 24px;
  font-weight: bold;
`;

const SignIn: React.FC = () => {
  const { login, submitTFA, loading, error } = useAuth();
  const [showTFA, setShowTFA] = useState(false);
  const [tfaToken, setTfaToken] = useState("");

  const onFinish = async (values: SignInCredentials) => {
    if (!showTFA) {
      const result = await login(values);
      if (result?.token) {
        setTfaToken(result.token);
        setShowTFA(true);
      }
    }
  };

  const onTFASubmit = async (values: { tfa_code: string }) => {
    const tfaRequest: TFARequest = {
      tfa_code: values.tfa_code,
      tfa_token: tfaToken,
      remember_me: true,
    };
    await submitTFA(tfaRequest);
  };

  return (
    <Container>
      <StyledCard>
        <Logo>Vetchi Employer Portal</Logo>
        {error && (
          <Alert
            message={error}
            type="error"
            showIcon
            style={{ marginBottom: 24 }}
          />
        )}

        {!showTFA ? (
          <>
            <Title level={3} style={{ textAlign: "center", marginBottom: 24 }}>
              Sign In
            </Title>
            <Form name="signin" onFinish={onFinish} layout="vertical">
              <Form.Item
                name="domain"
                rules={[
                  { required: true, message: "Please input your domain!" },
                ]}
              >
                <Input
                  prefix={<GlobalOutlined />}
                  placeholder="Domain"
                  size="large"
                />
              </Form.Item>
              <Form.Item
                name="email"
                rules={[
                  { required: true, message: "Please input your email!" },
                  { type: "email", message: "Please enter a valid email!" },
                ]}
              >
                <Input
                  prefix={<UserOutlined />}
                  placeholder="Email"
                  size="large"
                />
              </Form.Item>
              <Form.Item
                name="password"
                rules={[
                  { required: true, message: "Please input your password!" },
                ]}
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
          </>
        ) : (
          <>
            <Title level={3} style={{ textAlign: "center", marginBottom: 24 }}>
              Two-Factor Authentication
            </Title>
            <Form name="tfa" onFinish={onTFASubmit} layout="vertical">
              <Form.Item
                name="tfa_code"
                rules={[
                  { required: true, message: "Please input your TFA code!" },
                ]}
              >
                <Input placeholder="Enter TFA Code" size="large" />
              </Form.Item>
              <Form.Item>
                <Button
                  type="primary"
                  htmlType="submit"
                  size="large"
                  block
                  loading={loading}
                >
                  Verify
                </Button>
              </Form.Item>
            </Form>
          </>
        )}
      </StyledCard>
    </Container>
  );
};

export default SignIn;
