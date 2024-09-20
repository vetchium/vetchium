import { PlusCircleTwoTone } from "@ant-design/icons";
import { Button, Flex, Table } from "antd";
import { useNavigate } from "react-router-dom";
import { tableStyle } from "../Styles";
import t from "../i18n/i18n";

function Openings() {
  const navigate = useNavigate();

  const dataSource = [
    {
      key: "JAN14-1",
      id: "JAN14-1",
      status: "DRAFT",
      title: "Software Engineer",
      department: "Engineering",
      hiringManagerName: "John Doe",
      hiringManagerEmail: "john.doe@example.com",
      recruiterName: "Mother Theresa",
      recruiterEmail: "mother.theresa@example.com",
      unfilledPositions: 3,
      filledPositions: 2,
      editLink: <a href="#">Edit</a>,
      applicationsLink: (
        <a
          href={`/openings/JAN14-1/applications`}
          onClick={(e) => {
            e.preventDefault();
            navigate(`/openings/JAN14-1/applications`);
          }}
        >
          Applications
        </a>
      ),
    },
    {
      key: "JAN14-2",
      id: "JAN14-2",
      status: "ACTIVE",
      title: "Product Manager",
      department: "Product",
      hiringManagerName: "Jane Smith",
      hiringManagerEmail: "jane.smith@example.com",
      recruiterName: "Diana Prince",
      recruiterEmail: "diana.prince@example.com",
      unfilledPositions: 1,
      filledPositions: 4,
      editLink: <a href="#">Edit</a>,
      applicationsLink: (
        <a
          href={`/openings/JAN14-2/applications`}
          onClick={(e) => {
            e.preventDefault();
            navigate(`/openings/JAN14-2/applications`);
          }}
        >
          Applications
        </a>
      ),
    },
    {
      key: "FEB14-1",
      id: "FEB14-1",
      status: "CLOSED",
      title: "Sales Representative",
      department: "Sales",
      hiringManagerName: "Bob Johnson",
      hiringManagerEmail: "bob.johnson@example.com",
      recruiterName: "Mother Theresa",
      recruiterEmail: "mother.theresa@example.com",
      unfilledPositions: 0,
      filledPositions: 5,
      editLink: <a href="#">Edit</a>,
      applicationsLink: (
        <a
          href={`/openings/FEB14-1/applications`}
          onClick={(e) => {
            e.preventDefault();
            navigate(`/openings/FEB14-1/applications`);
          }}
        >
          Applications
        </a>
      ),
    },
  ];

  const columns = [
    { title: "ID", dataIndex: "id", key: "id" },
    { title: "Status", dataIndex: "status", key: "status" },
    { title: "Title", dataIndex: "title", key: "title" },
    { title: "Department", dataIndex: "department", key: "department" },
    {
      title: "Hiring Manager",
      dataIndex: "hiringManagerName",
      key: "hiringManagerName",
      render: (hiringManagerName: string, record: any) =>
        `${hiringManagerName} (${record.hiringManagerEmail})`,
      filters: [
        { text: "Jane Smith", value: "Jane Smith" },
        { text: "Bob Johnson", value: "Bob Johnson" },
      ],
      onFilter: (value: string, record: any) =>
        record.hiringManagerName.includes(value),
    },
    {
      title: "Recruiter",
      dataIndex: "recruiterName",
      key: "recruiterName",
      render: (recruiterName: string, record: any) =>
        `${recruiterName} (${record.recruiterEmail})`,
      filters: [
        { text: "Diana Prince", value: "Diana Prince" },
        { text: "Mother Theresa", value: "Mother Theresa" },
      ],
      onFilter: (value: string, record: any) =>
        record.recruiterName.includes(value),
    },
    {
      title: "Unfilled Positions",
      dataIndex: "unfilledPositions",
      key: "unfilledPositions",
    },
    {
      title: "Filled Positions",
      dataIndex: "filledPositions",
      key: "filledPositions",
    },
    { title: "Edit", dataIndex: "editLink", key: "editLink" },
    {
      title: "Applications",
      dataIndex: "applicationsLink",
      key: "applicationsLink",
    },
  ];

  return (
    <Flex wrap vertical>
      <Button
        type="primary"
        icon={<PlusCircleTwoTone />}
        onClick={() => navigate("/create-opening")}
        style={{ marginTop: "1rem", marginLeft: "2rem", width: "fit-content" }}
      >
        {t("openings.create_opening")}
      </Button>

      <Table
        dataSource={dataSource}
        columns={columns as any}
        style={tableStyle}
        scroll={{ x: true }}
      />
    </Flex>
  );
}

export default Openings;
