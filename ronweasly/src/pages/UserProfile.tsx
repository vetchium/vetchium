import { Button, Flex, Typography, Image, Card, Badge } from "antd";
import t from "../i18n/i18n";

const profileData = {
  name: "John Doe",
  photo: "https://placehold.co/120x120",
  description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
  work_history: [
    {
      company_name: "Example Private Limited",
      company_logo: "https://placehold.co/80x80",
      job_title: "Principal Software Engineer",
      start_date: "2022-01-01",
      end_date: "2023-01-01",
      still_employed: true,
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
      verified_colleagues: 102,
    },
    {
      company_name: "Example Private Limited",
      company_logo: "https://placehold.co/80x80",
      job_title: "Staff Engineer",
      start_date: "2021-01-01",
      end_date: "2022-01-01",
      still_employed: true,
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
      verified_colleagues: 10,
    },
    {
      company_name: "Example Private Limited",
      company_logo: "https://placehold.co/80x80",
      job_title: "Software Engineer",
      start_date: "2020-01-01",
      end_date: "2021-01-01",
      still_employed: true,
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
    },
  ],
};

function UserProfile() {
  return (
    <Flex gap="large" style={{ width: "100%", margin: "20px" }}>
      <Flex vertical>
        <Typography.Title
          level={2}
          style={{ fontSize: "24px", fontWeight: "bold" }}
        >
          {profileData.name}
        </Typography.Title>
        <Typography.Paragraph>{profileData.description}</Typography.Paragraph>
        <Flex vertical style={{ flex: 1, marginTop: "20px" }}>
          <Typography.Title level={3}>Work Experience</Typography.Title>
          {profileData.work_history.map((work) => (
            <Card>
              <Flex>
                <Image
                  src={work.company_logo}
                  alt="Company Logo"
                  style={{ borderRadius: "50%" }}
                />

                <Flex
                  key={work.company_name}
                  vertical
                  style={{ margin: "10px 0" }}
                >
                  <Typography.Title level={4}>
                    {work.job_title}
                  </Typography.Title>
                  <Typography.Text>{work.company_name}</Typography.Text>
                  <Typography.Text>
                    {work.start_date} - {work.end_date}
                  </Typography.Text>
                  <Typography.Text>{work.description}</Typography.Text>
                  {work.verified_colleagues ? (
                    <Flex gap="small">
                      <Badge
                        count={work.verified_colleagues}
                        style={{ backgroundColor: "#52c41a" }}
                      />
                      {t("user_profile.verified_colleagues")}
                    </Flex>
                  ) : (
                    ""
                  )}
                </Flex>
              </Flex>
            </Card>
          ))}
        </Flex>
      </Flex>
      <Flex vertical gap="small" style={{ marginLeft: "20%" }}>
        <Image
          src={profileData.photo}
          alt="Profile Photo"
          style={{ borderRadius: "50%" }}
        />
        <Flex vertical gap="small">
          <Button type="primary">Add as a colleague</Button>
          <Button>Follow posts</Button>
        </Flex>
      </Flex>
    </Flex>
  );
}

export default UserProfile;
