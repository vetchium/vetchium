import Typography from "@mui/material/Typography";
import Grid from "@mui/material/Grid";
import Paper from "@mui/material/Paper";
import Box from "@mui/material/Box";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";

export default function DashboardPage() {
  return (
    <AuthenticatedLayout>
      <Box sx={{ flexGrow: 1 }}>
        <Typography variant="h4" gutterBottom>
          Dashboard
        </Typography>
        <Grid container spacing={3}>
          <Grid item xs={12} md={6} lg={4}>
            <Paper
              sx={{
                p: 3,
                display: "flex",
                flexDirection: "column",
                height: 240,
              }}
            >
              <Typography variant="h6" gutterBottom>
                Recent Applications
              </Typography>
              {/* Add content here */}
            </Paper>
          </Grid>
          <Grid item xs={12} md={6} lg={4}>
            <Paper
              sx={{
                p: 3,
                display: "flex",
                flexDirection: "column",
                height: 240,
              }}
            >
              <Typography variant="h6" gutterBottom>
                Active Openings
              </Typography>
              {/* Add content here */}
            </Paper>
          </Grid>
          <Grid item xs={12} md={6} lg={4}>
            <Paper
              sx={{
                p: 3,
                display: "flex",
                flexDirection: "column",
                height: 240,
              }}
            >
              <Typography variant="h6" gutterBottom>
                Upcoming Interviews
              </Typography>
              {/* Add content here */}
            </Paper>
          </Grid>
        </Grid>
      </Box>
    </AuthenticatedLayout>
  );
}
