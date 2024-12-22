import React, { useState, useEffect } from "react";
import { Row, Col, Card, Statistic, List, Typography, Spin } from "antd";
import {
  FileSearchOutlined,
  FileTextOutlined,
  UserOutlined,
  CheckCircleOutlined,
} from "@ant-design/icons";
import axios from "axios";
import { useAuth } from "@/hooks/useAuth";
import { Application } from "@/types/application";
import { Candidacy } from "@/types/candidacy";

const { Title } = Typography;

interface DashboardStats {
  total_applications: number;
  active_candidacies: number;
  shortlisted_applications: number;
  upcoming_interviews: number;
}

const Dashboard: React.FC = () => {
  const { user } = useAuth();
  const [stats, setStats] = useState<DashboardStats>({
    total_applications: 0,
    active_candidacies: 0,
    shortlisted_applications: 0,
    upcoming_interviews: 0,
  });
  const [recentApplications, setRecentApplications] = useState<Application[]>(
    []
  );
  const [recentCandidacies, setRecentCandidacies] = useState<Candidacy[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        const [statsRes, applicationsRes, candidaciesRes] = await Promise.all([
          axios.get<DashboardStats>("/api/hub/dashboard/stats"),
          axios.get<{ items: Application[] }>("/api/hub/applications", {
            params: { limit: 5 },
          }),
          axios.get<{ items: Candidacy[] }>("/api/hub/candidacies", {
            params: { limit: 5 },
          }),
        ]);

        setStats(statsRes.data);
        setRecentApplications(applicationsRes.data.items);
        setRecentCandidacies(candidaciesRes.data.items);
      } catch (error) {
        console.error("Failed to fetch dashboard data:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  if (loading) {
    return (
      <div style={{ textAlign: "center", padding: "50px" }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div>
      <Title level={2}>Welcome back, {user?.name}!</Title>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Total Applications"
              value={stats.total_applications}
              prefix={<FileTextOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Active Candidacies"
              value={stats.active_candidacies}
              prefix={<UserOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Shortlisted"
              value={stats.shortlisted_applications}
              prefix={<FileSearchOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Upcoming Interviews"
              value={stats.upcoming_interviews}
              prefix={<CheckCircleOutlined />}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col xs={24} lg={12}>
          <Card title="Recent Applications">
            <List
              dataSource={recentApplications}
              renderItem={(item) => (
                <List.Item>
                  <List.Item.Meta
                    title={item.opening.title}
                    description={`${item.opening.employer_name} - ${item.state}`}
                  />
                </List.Item>
              )}
              locale={{ emptyText: "No recent applications" }}
            />
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card title="Active Candidacies">
            <List
              dataSource={recentCandidacies}
              renderItem={(item) => (
                <List.Item>
                  <List.Item.Meta
                    title={item.application.opening.title}
                    description={`${item.application.opening.employer_name} - ${item.state}`}
                  />
                </List.Item>
              )}
              locale={{ emptyText: "No active candidacies" }}
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
