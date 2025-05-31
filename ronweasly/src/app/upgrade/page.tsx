"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { useAuth } from "@/hooks/useAuth";
import { useTranslation } from "@/hooks/useTranslation";
import CheckIcon from "@mui/icons-material/Check";
import StarIcon from "@mui/icons-material/Star";
import {
  Box,
  Button,
  Card,
  CardContent,
  Container,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Typography,
} from "@mui/material";

export default function UpgradePage() {
  const { t } = useTranslation();
  useAuth(); // Check authentication and redirect if not authenticated

  const paidFeatures = [
    "Unlimited post length (up to 4096 characters)",
    "Create new tags when posting",
    "Change your handle",
    "Upload profile pictures",
    "No advertisements",
    "Support open source development",
  ];

  return (
    <AuthenticatedLayout>
      <Container maxWidth="md" sx={{ mt: 4, mb: 4 }}>
        <Box sx={{ textAlign: "center", mb: 4 }}>
          <Typography variant="h3" component="h1" gutterBottom>
            Upgrade to Paid Tier
          </Typography>
          <Typography variant="h6" color="text.secondary">
            Unlock all features and support the platform
          </Typography>
        </Box>

        <Card sx={{ maxWidth: 600, mx: "auto", p: 3 }}>
          <CardContent>
            <Box sx={{ textAlign: "center", mb: 3 }}>
              <StarIcon sx={{ fontSize: 48, color: "primary.main", mb: 2 }} />
              <Typography variant="h4" component="h2" gutterBottom>
                Paid Tier
              </Typography>
              <Typography variant="h3" color="primary" sx={{ mb: 2 }}>
                $99/year
              </Typography>
            </Box>

            <List>
              {paidFeatures.map((feature, index) => (
                <ListItem key={index}>
                  <ListItemIcon>
                    <CheckIcon color="primary" />
                  </ListItemIcon>
                  <ListItemText primary={feature} />
                </ListItem>
              ))}
            </List>

            <Box sx={{ textAlign: "center", mt: 4 }}>
              <Button
                variant="contained"
                size="large"
                sx={{ px: 4, py: 2 }}
                onClick={() => {
                  // TODO: Implement actual upgrade functionality
                  alert("Upgrade functionality coming soon!");
                }}
              >
                Upgrade Now
              </Button>
            </Box>

            <Typography
              variant="body2"
              color="text.secondary"
              sx={{ textAlign: "center", mt: 3 }}
            >
              This is a placeholder page. Actual payment processing will be
              implemented soon.
            </Typography>
          </CardContent>
        </Card>
      </Container>
    </AuthenticatedLayout>
  );
}
