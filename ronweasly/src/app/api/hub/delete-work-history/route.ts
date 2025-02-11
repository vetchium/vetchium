import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const response = await fetch(
      `${process.env.API_BASE_URL}/hub/delete-work-history`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: request.headers.get("Authorization") || "",
        },
        body: JSON.stringify(body),
      }
    );

    if (!response.ok) {
      throw new Error("Failed to delete work history");
    }

    return NextResponse.json({ success: true });
  } catch (error) {
    console.error("Error in delete-work-history:", error);
    return NextResponse.json(
      { error: "Failed to delete work history" },
      { status: 500 }
    );
  }
}
