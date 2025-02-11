import { NextResponse } from "next/server";

export async function PUT(request: Request) {
  try {
    const body = await request.json();
    const response = await fetch(
      `${process.env.API_BASE_URL}/hub/update-work-history`,
      {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: request.headers.get("Authorization") || "",
        },
        body: JSON.stringify(body),
      }
    );

    if (!response.ok) {
      throw new Error("Failed to update work history");
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error("Error in update-work-history:", error);
    return NextResponse.json(
      { error: "Failed to update work history" },
      { status: 500 }
    );
  }
}
