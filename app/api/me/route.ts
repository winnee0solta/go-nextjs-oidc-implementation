import { NextRequest, NextResponse } from "next/server";

export async function GET(req: NextRequest) {
  const cookieHeader = req.headers.get("cookie") || "";

  const res = await fetch("http://localhost:8080/me", {
    headers: {
      Cookie: cookieHeader,
    },
    credentials: "include",
  });

  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}
