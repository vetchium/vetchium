/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",
  reactStrictMode: true,
  // Make API_ENDPOINT available at runtime
  env: {
    API_ENDPOINT: process.env.API_ENDPOINT || "http://localhost:8081",
  },
  async rewrites() {
    // Use the runtime env variable
    return [
      {
        source: "/api/:path*",
        destination: `${
          process.env.API_ENDPOINT || "http://localhost:8081"
        }/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;
