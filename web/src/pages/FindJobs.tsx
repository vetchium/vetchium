import {
  Flex,
  Typography,
  Checkbox,
  InputNumber,
  Select,
  List,
  Avatar,
  Card,
  Slider,
  Tag,
  AutoComplete,
} from "antd";
import { useState } from "react";

// Generate 50 rows of data
const openingsData = Array.from({ length: 50 }, (_, i) => ({
  id: i + 1,
  title: [
    "Software Engineer",
    "Product Manager",
    "Data Scientist",
    "UX Designer",
    "Marketing Specialist",
  ][Math.floor(Math.random() * 5)],
  company: ["TechCorp", "DataInc", "DesignHub", "MarketPro", "AITech"][
    Math.floor(Math.random() * 5)
  ],
  location: ["San Francisco", "New York", "London", "Berlin", "Tokyo"][
    Math.floor(Math.random() * 5)
  ],
  jobType: ["full_time", "part_time", "contract", "internship"][
    Math.floor(Math.random() * 4)
  ],
  salary: {
    min: Math.floor(Math.random() * 50000) + 30000,
    max: Math.floor(Math.random() * 100000) + 80000,
    currency: ["USD", "EUR", "GBP", "JPY"][Math.floor(Math.random() * 4)],
  },
  experience: {
    min: Math.floor(Math.random() * 5),
    max: Math.floor(Math.random() * 10) + 5,
  },
  educationalQualification: ["bachelors", "masters", "phd", "unspecified"][
    Math.floor(Math.random() * 4)
  ],
  remoteLocations: ["USA", "Europe", "Asia", "Global"][
    Math.floor(Math.random() * 4)
  ],
}));

const locationOptions = [
  { value: "USA - California - San Francisco" },
  { value: "USA - New York - New York City" },
  { value: "UK - England - London" },
  { value: "Germany - Berlin - Berlin" },
  { value: "Japan - Tokyo - Tokyo" },
  // Add more options as needed
];

function FindJobs() {
  const [filteredJobs, setFilteredJobs] = useState(openingsData);
  const [selectedLocations, setSelectedLocations] = useState<string[]>([]);

  // Filter function (placeholder)
  const filterJobs = () => {
    // Implement filtering logic here
    setFilteredJobs(openingsData);
  };

  const handleLocationSelect = (value: string) => {
    setSelectedLocations((prev) => [...prev, value]);
    filterJobs();
  };

  const handleLocationCheckboxChange = (checkedValues: string[]) => {
    setSelectedLocations(checkedValues);
    filterJobs();
  };

  return (
    <Flex vertical gap="large">
      <Typography.Title level={2}>Find Jobs</Typography.Title>
      <Flex>
        {/* Left column: Filters */}
        <Flex vertical style={{ width: "30%", marginRight: "20px" }}>
          <Card title="Filter Jobs">
            <Flex vertical gap="middle">
              <Typography.Text strong>Job Type</Typography.Text>
              <Checkbox.Group
                options={["Full-time", "Part-time", "Contract", "Internship"]}
                onChange={filterJobs}
              />

              <Typography.Text strong>Locations</Typography.Text>
              <AutoComplete
                style={{ width: "100%" }}
                options={locationOptions}
                placeholder="Type to search locations"
                onSelect={handleLocationSelect}
              />
              <Checkbox.Group
                options={selectedLocations.map((loc) => ({
                  label: loc,
                  value: loc,
                }))}
                value={selectedLocations}
                onChange={handleLocationCheckboxChange}
              />

              <Typography.Text strong>Salary Range</Typography.Text>
              <Flex gap="small">
                <InputNumber
                  style={{ width: "50%" }}
                  placeholder="Min"
                  onChange={filterJobs}
                />
                <InputNumber
                  style={{ width: "50%" }}
                  placeholder="Max"
                  onChange={filterJobs}
                />
              </Flex>

              <Typography.Text strong>Experience (Years)</Typography.Text>
              <Slider
                range
                min={0}
                max={20}
                defaultValue={[0, 20]}
                onChange={filterJobs}
              />

              <Typography.Text strong>
                Educational Qualification
              </Typography.Text>
              <Checkbox.Group
                options={[
                  { value: "bachelors", label: "Bachelor's" },
                  { value: "masters", label: "Master's" },
                  { value: "phd", label: "PhD" },
                  { value: "unspecified", label: "Unspecified" },
                ]}
                onChange={filterJobs}
              />

              <Typography.Text strong>Remote Locations</Typography.Text>
              <Select
                mode="multiple"
                style={{ width: "100%" }}
                placeholder="Select remote locations"
                onChange={filterJobs}
                options={[
                  { value: "usa", label: "USA" },
                  { value: "europe", label: "Europe" },
                  { value: "asia", label: "Asia" },
                  { value: "global", label: "Global" },
                ]}
              />
            </Flex>
          </Card>
        </Flex>

        {/* Right column: Job listings */}
        <Flex vertical style={{ width: "70%" }}>
          <List
            itemLayout="vertical"
            size="large"
            pagination={{
              onChange: (page) => {
                console.log(page);
              },
              pageSize: 10,
            }}
            dataSource={filteredJobs}
            renderItem={(job) => (
              <List.Item
                key={job.id}
                extra={
                  <Avatar
                    size={64}
                    src={`https://xsgames.co/randomusers/avatar.php?g=pixel&key=${job.id}`}
                  />
                }
              >
                <List.Item.Meta
                  title={
                    <Typography.Link href={`/job/${job.id}`}>
                      {job.title}
                    </Typography.Link>
                  }
                  description={<Typography.Text>{job.company}</Typography.Text>}
                />
                <Flex vertical>
                  <Typography.Text>{job.location}</Typography.Text>
                  <Typography.Text>
                    <Tag color="blue">{job.jobType}</Tag>
                  </Typography.Text>
                  <Typography.Text>{`${job.salary.currency} ${job.salary.min} - ${job.salary.max}`}</Typography.Text>
                  <Typography.Text>{`${job.experience.min} - ${job.experience.max} years experience`}</Typography.Text>
                  <Typography.Text>{`Education: ${job.educationalQualification}`}</Typography.Text>
                  <Typography.Text>{`Remote: ${job.remoteLocations}`}</Typography.Text>
                </Flex>
              </List.Item>
            )}
          />
        </Flex>
      </Flex>
    </Flex>
  );
}

export default FindJobs;
