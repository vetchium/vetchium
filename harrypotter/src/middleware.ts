import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const PUBLIC_PATHS = ["/signin", "/tfa"];

export default async function middleware(request: NextRequest) {
  const publicPath = PUBLIC_PATHS.some((path) =>
    request.nextUrl.pathname.startsWith(path)
  );

  if (publicPath) {
    return NextResponse.next();
  }

  const token =
    request.cookies.get("sessionToken") ||
    request.headers.get("Authorization")?.split(" ")[1];

  if (!token) {
    return NextResponse.redirect(new URL("/signin", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!api|_next|.*\\..*).*)"],
};
