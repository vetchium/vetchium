import { Table } from "antd";

function Openings() {
  const dataSource = [
    {
      key: "1",
      id: "1",
      status: "DRAFT",
      title: "Software Engineer",
      department: "Engineering",
      hiringManager: "John Doe",
      unfilledPositions: 3,
      filledPositions: 2,
      createdAt: "2022-01-01",
      editLink: <a href="#">Edit</a>,
    },
    {
      key: "2",
      id: "2",
      status: "ACTIVE",
      title: "Product Manager",
      department: "Product",
      hiringManager: "Jane Smith",
      unfilledPositions: 1,
      filledPositions: 4,
      createdAt: "2022-01-02",
      editLink: <a href="#">Edit</a>,
    },
    {
      key: "3",
      id: "3",
      status: "CLOSED",
      title: "Sales Representative",
      department: "Sales",
      hiringManager: "Bob Johnson",
      unfilledPositions: 0,
      filledPositions: 5,
      createdAt: "2022-01-03",
      editLink: <a href="#">Edit</a>,
    },
  ];

  const columns = [
    { title: "ID", dataIndex: "id", key: "id" },
    { title: "Status", dataIndex: "status", key: "status" },
    { title: "Title", dataIndex: "title", key: "title" },
    { title: "Department", dataIndex: "department", key: "department" },
    {
      title: "Hiring Manager",
      dataIndex: "hiringManager",
      key: "hiringManager",
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
    { title: "Created At", dataIndex: "createdAt", key: "createdAt" },
    { title: "Edit", dataIndex: "editLink", key: "editLink" },
  ];

  return <Table dataSource={dataSource} columns={columns} />;
}

export default Openings;
