import { Flex, Image, List, Typography } from "antd";

const data = [
  {
    company_name: "Anywhere Private Limited",
    job_title: "Principal Engineer",
    start_date: "Jan 2024",
    end_date: "Present",
    logo: "https://placehold.co/80x80",
  },
  {
    company_name: "Anywhere Private Limited",
    job_title: "Staff Engineer",
    start_date: "Jan 2023",
    end_date: "Dec 2023",
    logo: "https://placehold.co/80x80",
  },
  {
    company_name: "Something Private Limited",
    job_title: "Software Architect",
    start_date: "Jan 2022",
    end_date: "Dec 2022",
    logo: "https://placehold.co/80x80",
  },
  {
    company_name: "Whatever Private Limited",
    job_title: "Senior Software Engineer",
    start_date: "Jan 2021",
    end_date: "Dec 2021",
    logo: "https://placehold.co/80x80",
  },
  {
    company_name: "Example Inc.",
    job_title: "Software Engineer",
    start_date: "Jan 2020",
    end_date: "Dec 2020",
    logo: "https://placehold.co/80x80",
  },
];

function WorkHistory() {
  return (
    <Flex vertical gap="large">
      <Typography.Title level={2}>Work History</Typography.Title>
      <List
        itemLayout="vertical"
        dataSource={data}
        renderItem={(item: any) => (
          <List.Item>
            <Flex gap="large">
              <Image src={item.logo} />
              <Flex vertical>
                <Typography.Title level={4}>{item.job_title}</Typography.Title>
                <Typography.Text>{item.company_name}</Typography.Text>
                <Typography.Text>
                  {item.start_date} - {item.end_date}
                </Typography.Text>
              </Flex>
            </Flex>
          </List.Item>
        )}
      />
    </Flex>
  );
}

export default WorkHistory;
