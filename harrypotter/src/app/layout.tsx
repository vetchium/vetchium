import ThemeRegistry from "@/components/ThemeRegistry";

export const metadata = {
  title: "Vetchi",
  description: "Vetchi Employer Portal",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <ThemeRegistry>{children}</ThemeRegistry>
      </body>
    </html>
  );
}
