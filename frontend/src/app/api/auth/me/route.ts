import { NextResponse } from "next/server";
import { cookies } from "next/headers";

export async function GET(request: Request) {
  let token: string | undefined;

  const authHeader = request.headers.get("Authorization");
  if (authHeader && authHeader.startsWith("Bearer ")) {
    token = authHeader.replace("Bearer ", "").trim();
  }

  if (!token) {
    const cookieStore = await cookies();
    token = cookieStore.get("gh_session")?.value;
  }

  if (!token) {
    token = process.env.GITHUB_TOKEN_CLASSIC;
  }

  if (!token) {
    return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
  }

  try {
    const res = await fetch("https://api.github.com/user", {
      headers: {
        Accept: "application/vnd.github+json",
        Authorization: `Bearer ${token}`,
        "X-GitHub-Api-Version": "2022-11-28",
      },
    });

    if (!res.ok) {
      return NextResponse.json(
        { error: `GitHub API returned ${res.status}` },
        { status: res.status }
      );
    }

    const data = await res.json();
    return NextResponse.json({
      login: data.login,
      id: data.id,
      avatar_url: data.avatar_url,
      name: data.name,
      email: data.email,
    });
  } catch (err: any) {
    return NextResponse.json({ error: err.message }, { status: 500 });
  }
}
