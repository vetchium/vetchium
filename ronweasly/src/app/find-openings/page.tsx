import { useState } from "react";
import Box from "@mui/material/Box";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";
import Paper from "@mui/material/Paper";
import SearchIcon from "@mui/icons-material/Search";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";

export default function FindOpeningsPage() {
  const [searchQuery, setSearchQuery] = useState("");

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: Implement search functionality
    console.log("Searching for:", searchQuery);
  };

  return (
    <AuthenticatedLayout>
      <Box sx={{ maxWidth: 800, mx: "auto", mt: 4 }}>
        <Typography variant="h4" gutterBottom align="center">
          Find Openings
        </Typography>
        <Paper
          component="form"
          onSubmit={handleSearch}
          sx={{
            p: 4,
            mt: 4,
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
          }}
        >
          <Typography
            variant="body1"
            color="text.secondary"
            gutterBottom
            align="center"
          >
            Search for job openings across all locations
          </Typography>
          <Box sx={{ width: "100%", mt: 2 }}>
            <TextField
              fullWidth
              variant="outlined"
              placeholder="Search for job titles, skills, or keywords"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              InputProps={{
                endAdornment: (
                  <Button
                    variant="contained"
                    sx={{ ml: 1 }}
                    type="submit"
                    startIcon={<SearchIcon />}
                  >
                    Search
                  </Button>
                ),
              }}
            />
          </Box>
        </Paper>
        {/* Search results will be displayed here */}
      </Box>
    </AuthenticatedLayout>
  );
}
