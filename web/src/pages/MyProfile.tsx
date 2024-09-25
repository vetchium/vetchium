import { Card, Flex, Image, Typography } from "antd";
import WorkHistory from "../components/WorkHistory";
import t from "../i18n/i18n";
import AutoBiography from "../components/AutoBiography";

function MyProfile() {
  return (
    <Flex vertical gap="large" style={{ margin: "2rem" }}>
      <Image
        src={"https://placehold.co/120x120"}
        style={{ borderRadius: "50%", width: "120px", height: "120px" }}
      />
      <Flex vertical gap="small" style={{ margin: "2rem" }}>
        <AutoBiography />
        <WorkHistory />

        <Card>
          <Typography.Title level={3}>
            {t("myprofile.education")}
          </Typography.Title>
          <Typography.Text>
            Lorem ipsum dolor sit amet, consectetur adipiscing elit.
          </Typography.Text>
        </Card>

        <Card>
          <Typography.Title level={3}>
            {t("myprofile.certifications")}
          </Typography.Title>
          <Typography.Text>
            Lorem ipsum dolor sit amet, consectetur adipiscing elit.
          </Typography.Text>
        </Card>

        <Card>
          <Typography.Title level={3}>
            {t("myprofile.patents_publications")}
          </Typography.Title>
          <Typography.Text>
            Lorem ipsum dolor sit amet, consectetur adipiscing elit.
          </Typography.Text>
        </Card>

        <Card>
          <Typography.Title level={3}>{t("myprofile.skills")}</Typography.Title>
          <Typography.Text>
            Lorem ipsum dolor sit amet, consectetur adipiscing elit.
          </Typography.Text>
        </Card>

        <Card>
          <Typography.Title level={3}>
            {t("myprofile.languages")}
          </Typography.Title>
          <Typography.Text>
            Lorem ipsum dolor sit amet, consectetur adipiscing elit.
          </Typography.Text>
        </Card>
      </Flex>
    </Flex>
  );
}

export default MyProfile;
