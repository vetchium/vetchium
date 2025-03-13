/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",
  reactStrictMode: true,
  async rewrites() {
    // Get API endpoint from env var
    const apiEndpoint = process.env.API_ENDPOINT || "http://VETCHI_MISSED_URL";
    return [
      {
        source: "/api/:path*",
        destination: `${apiEndpoint}/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;
