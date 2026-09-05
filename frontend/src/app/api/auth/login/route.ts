import { NextResponse } from "next/server";

export async function GET(request: Request) {
  const clientId = process.env.GITHUB_CLIENT_ID;
  if (!clientId) {
    return NextResponse.json(
      { error: "GITHUB_CLIENT_ID not configured in environment." },
      { status: 400 }
    );
  }

  const url = new URL(request.url);
  const rawAppUrl = process.env.APP_URL;
  const baseUrl = rawAppUrl && !rawAppUrl.includes("localhost")
    ? rawAppUrl.replace(/\/$/, "")
    : url.origin;
  const redirectUri = `${baseUrl}/api/auth/callback`;
  const state = Math.random().toString(36).substring(2, 15);
  const scope = "repo,admin:repo_hook";

  const githubAuthUrl = `https://github.com/login/oauth/authorize?client_id=${clientId}&redirect_uri=${encodeURIComponent(
    redirectUri
  )}&scope=${encodeURIComponent(scope)}&state=${state}`;

  return NextResponse.redirect(githubAuthUrl);
}
