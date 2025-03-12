/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",
  reactStrictMode: true,
  async rewrites() {
    // Get API endpoint from env var with localhost fallback for development
    const apiEndpoint = process.env.API_ENDPOINT || "http://localhost:8081";
    return [
      {
        source: "/api/:path*",
        destination: `${apiEndpoint}/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;
