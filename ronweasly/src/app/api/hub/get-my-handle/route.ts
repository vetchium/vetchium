import { NextResponse } from "next/server";

export async function GET(request: Request) {
  try {
    const response = await fetch(
      `${process.env.API_BASE_URL}/hub/get-my-handle`,
      {
        headers: {
          Authorization: request.headers.get("Authorization") || "",
        },
      }
    );

    if (!response.ok) {
      throw new Error("Failed to fetch user handle");
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error("Error in get-my-handle:", error);
    return NextResponse.json(
      { error: "Failed to fetch user handle" },
      { status: 500 }
    );
  }
}
