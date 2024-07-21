import { Button, Flex, List, Select, Typography } from "antd";
import t from "../i18n/i18n";

const { Text, Link } = Typography;

const data = [
  {
    name: "John Doe",
    currentPosition: "Software Engineer, Company A",
    vouchedBy: [
      { name: "Jane Smith", position: "Manager", company: "Company A" },
      { name: "Bob Johnson", position: "Team Lead", company: "Company X" },
    ],
    referredBy: "Sankar P <sankar@example.com>",
  },
  {
    name: "Alice Brown",
    currentPosition: "Product Manager, Company B",
    vouchedBy: [
      { name: "Charlie Davis", position: "Director", company: "Company B" },
      { name: "Eve White", position: "CEO", company: "Company Y" },
    ],
  },
  {
    name: "Michael Green",
    currentPosition: "Data Scientist, Company C",
    vouchedBy: [
      { name: "Fiona Black", position: "CTO", company: "Company C" },
      {
        name: "George Blue",
        position: "Lead Data Scientist",
        company: "Company Z",
      },
    ],
  },
  {
    name: "Sankarasivasubramanian Pasupathilingam",
    currentPosition: "UX Designer, Company D",
    vouchedBy: [
      { name: "Hannah Red", position: "Head of Design", company: "Company D" },
      { name: "Ian Yellow", position: "Senior Designer", company: "Company W" },
    ],
  },
];

export default function Applications() {
  return (
    <List
      itemLayout="vertical"
      dataSource={data}
      renderItem={(item) => (
        <List.Item style={{ marginBottom: "36px" }}>
          <Flex vertical gap="large">
            <Flex justify="space-between">
              <Flex vertical>
                <Typography.Title level={3} style={{ margin: 0 }}>
                  {item.name}
                </Typography.Title>
                <Text>{item.currentPosition}</Text>
              </Flex>

              <Flex vertical gap="small">
                <Text strong>{t("applications.vouched_by")}</Text>
                {item.vouchedBy.map((vouched, index) => (
                  <Link href="#">
                    {vouched.name} ({vouched.position}, {vouched.company})
                  </Link>
                ))}
                {item.referredBy && (
                  <>
                    <Text strong>{t("applications.referred_by")}</Text>
                    <Text>{item.referredBy}</Text>
                  </>
                )}
              </Flex>

              <Link href="#">{t("applications.resume")}</Link>
              <Flex vertical gap="small">
                <Select
                  mode="multiple"
                  placeholder={t("applications.add_tags")}
                  options={[
                    {
                      value: "red",
                      label: (
                        <div
                          style={{
                            width: "20px",
                            height: "20px",
                            backgroundColor: "red",
                            borderRadius: "50%",
                            display: "inline-block",
                          }}
                        ></div>
                      ),
                    },
                    {
                      value: "blue",
                      label: (
                        <div
                          style={{
                            width: "20px",
                            height: "20px",
                            backgroundColor: "blue",
                            borderRadius: "50%",
                            display: "inline-block",
                          }}
                        ></div>
                      ),
                    },
                    {
                      value: "green",
                      label: (
                        <div
                          style={{
                            width: "20px",
                            height: "20px",
                            backgroundColor: "green",
                            borderRadius: "50%",
                            display: "inline-block",
                          }}
                        ></div>
                      ),
                    },
                    {
                      value: "purple",
                      label: (
                        <div
                          style={{
                            width: "20px",
                            height: "20px",
                            backgroundColor: "purple",
                            borderRadius: "50%",
                            display: "inline-block",
                          }}
                        ></div>
                      ),
                    },
                  ]}
                  getPopupContainer={(triggerNode) => triggerNode.parentElement}
                />
                <Button type="primary">
                  {t("applications.schedule_interview")}
                </Button>
              </Flex>
            </Flex>
            <Flex>
              <Button danger>{t("applications.reject")}</Button>
            </Flex>
          </Flex>
        </List.Item>
      )}
    />
  );
}
