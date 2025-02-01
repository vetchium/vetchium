/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  async rewrites() {
    return [
      {
        source: "/api/:path*", // Match all requests to /api
        destination: "http://localhost:8081/:path*", // Proxy to your API server
      },
    ];
  },
};

module.exports = nextConfig;
