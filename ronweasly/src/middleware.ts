import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";

const publicPaths = ["/login", "/tfa"];
const signupPathRegex = /^\/signup-hubuser\/.+/;

export function middleware(request: NextRequest) {
  const sessionToken = request.cookies.get("session_token")?.value;
  const tfaToken = request.cookies.get("tfa_token")?.value;
  const pathname = request.nextUrl.pathname;

  // Allow signup paths
  if (signupPathRegex.test(pathname)) {
    return NextResponse.next();
  }

  // Allow public paths
  if (publicPaths.includes(pathname)) {
    // If user is already authenticated, redirect to home
    if (sessionToken) {
      return NextResponse.redirect(new URL("/", request.url));
    }
    return NextResponse.next();
  }

  // Check for TFA page
  if (pathname === "/tfa") {
    if (!tfaToken) {
      return NextResponse.redirect(new URL("/login", request.url));
    }
    return NextResponse.next();
  }

  // Protected routes
  if (!sessionToken) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  return NextResponse.next();
}

export const config = {
  // Match all request paths except for the ones starting with:
  // - api (API routes)
  // - _next/static (static files)
  // - _next/image (image optimization files)
  // - any path containing a dot (likely static files: .png, .ico, .webp, etc.)
  matcher: ["/((?!api|_next/static|_next/image|.*\\..*).*) "],
};
