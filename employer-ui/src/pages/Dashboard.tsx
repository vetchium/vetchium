import React from "react";
import { Row, Col, Card, Statistic } from "antd";
import {
  UserOutlined,
  FileTextOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
} from "@ant-design/icons";
import styled from "styled-components";

const StyledCard = styled(Card)`
  margin-bottom: 24px;
`;

const Dashboard: React.FC = () => {
  return (
    <div>
      <h1>Dashboard</h1>
      <Row gutter={24}>
        <Col xs={24} sm={12} lg={6}>
          <StyledCard>
            <Statistic
              title="Active Openings"
              value={5}
              prefix={<FileTextOutlined />}
            />
          </StyledCard>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StyledCard>
            <Statistic
              title="Total Applications"
              value={25}
              prefix={<UserOutlined />}
            />
          </StyledCard>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StyledCard>
            <Statistic
              title="Shortlisted"
              value={10}
              prefix={<CheckCircleOutlined />}
            />
          </StyledCard>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StyledCard>
            <Statistic
              title="Pending Reviews"
              value={8}
              prefix={<ClockCircleOutlined />}
            />
          </StyledCard>
        </Col>
      </Row>

      <Row gutter={24}>
        <Col xs={24} lg={12}>
          <StyledCard title="Recent Applications">
            {/* Add recent applications list component here */}
          </StyledCard>
        </Col>
        <Col xs={24} lg={12}>
          <StyledCard title="Upcoming Interviews">
            {/* Add upcoming interviews list component here */}
          </StyledCard>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
