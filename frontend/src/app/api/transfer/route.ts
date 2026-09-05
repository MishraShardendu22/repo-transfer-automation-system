import { NextResponse } from "next/server";
import { cookies } from "next/headers";

export async function POST(request: Request) {
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
    return NextResponse.json(
      { error: "Unauthorized: No GitHub token provided." },
      { status: 401 }
    );
  }

  try {
    const body = await request.json();
    const { target_user, repositories } = body;

    if (!target_user || !Array.isArray(repositories) || repositories.length === 0) {
      return NextResponse.json(
        { error: "target_user and at least one repository are required." },
        { status: 400 }
      );
    }

    // Determine current user
    const userRes = await fetch("https://api.github.com/user", {
      headers: {
        Accept: "application/vnd.github+json",
        Authorization: `Bearer ${token}`,
        "X-GitHub-Api-Version": "2022-11-28",
      },
    });

    if (!userRes.ok) {
      return NextResponse.json(
        { error: "Could not authenticate user with GitHub." },
        { status: 401 }
      );
    }

    const currentUser = await userRes.json();
    const originalUser = currentUser.login;

    const results = await Promise.all(
      repositories.map(async (repoName: string) => {
        try {
          const res = await fetch(
            `https://api.github.com/repos/${originalUser}/${repoName}/transfer`,
            {
              method: "POST",
              headers: {
                Accept: "application/vnd.github+json",
                Authorization: `Bearer ${token}`,
                "X-GitHub-Api-Version": "2022-11-28",
                "Content-Type": "application/json",
              },
              body: JSON.stringify({ new_owner: target_user }),
            }
          );

          if (res.status === 202) {
            return {
              repo: repoName,
              success: true,
              status_code: 202,
              message: `Transfer initiated successfully to ${target_user}.`,
            };
          }

          const errData = await res.json().catch(() => ({}));
          return {
            repo: repoName,
            success: false,
            status_code: res.status,
            message: errData.message || `GitHub error (status ${res.status})`,
          };
        } catch (err: any) {
          return {
            repo: repoName,
            success: false,
            status_code: 500,
            message: err.message,
          };
        }
      })
    );

    return NextResponse.json(results);
  } catch (err: any) {
    return NextResponse.json({ error: err.message }, { status: 500 });
  }
}
