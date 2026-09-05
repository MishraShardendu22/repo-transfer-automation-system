"use client";

import { useEffect, useState } from "react";
import { 
  ArrowRightLeft, 
  GitPullRequest, 
  RefreshCw, 
  Search, 
  CheckCircle2, 
  AlertCircle, 
  Terminal,
  LogOut,
  FolderGit2
} from "lucide-react";

interface User {
  login: string;
  avatar_url: string;
}

interface Repository {
  name: string;
  full_name: string;
  private: boolean;
  updated_at: string;
  description: string | null;
}

interface TransferResult {
  repo: string;
  success: boolean;
  status_code: number;
  message: string;
}

export default function Home() {
  const [user, setUser] = useState<User | null>(null);
  const [repos, setRepos] = useState<Repository[]>([]);
  const [search, setSearch] = useState("");
  const [selectedRepos, setSelectedRepos] = useState<string[]>([]);
  const [targetUser, setTargetUser] = useState("MishraShardendu22");
  const [patToken, setPatToken] = useState("");
  const [loading, setLoading] = useState(false);
  const [transferring, setTransferring] = useState(false);
  const [results, setResults] = useState<TransferResult[]>([]);
  const [logs, setLogs] = useState<{ text: string; type: "info" | "success" | "error"; time: string }[]>([
    { text: "System initialized. Authenticate or provide a token to begin.", type: "info", time: new Date().toLocaleTimeString() }
  ]);

  const addLog = (text: string, type: "info" | "success" | "error" = "info") => {
    setLogs((prev) => [...prev, { text, type, time: new Date().toLocaleTimeString() }]);
  };

  const fetchAuth = async () => {
    try {
      const res = await fetch("/api/auth/me");
      if (res.ok) {
        const data = await res.json();
        setUser(data);
        addLog(`Authenticated as ${data.login} via GitHub OAuth.`, "success");
        fetchRepos();
      } else {
        setUser(null);
      }
    } catch {
      setUser(null);
    }
  };

  const fetchRepos = async () => {
    setLoading(true);
    const headers: Record<string, string> = {};
    if (patToken.trim()) {
      headers["Authorization"] = `Bearer ${patToken.trim()}`;
    }

    try {
      const res = await fetch("/api/repos", { headers });
      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Failed to fetch repositories");
      }
      const data = await res.json();
      setRepos(data);
      addLog(`Fetched ${data.length} repositories from GitHub.`, "info");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Error fetching repositories";
      addLog(msg, "error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAuth();
  }, []);

  const filteredRepos = repos.filter((r) =>
    r.name.toLowerCase().includes(search.toLowerCase())
  );

  const toggleSelectAll = () => {
    if (selectedRepos.length === filteredRepos.length) {
      setSelectedRepos([]);
    } else {
      setSelectedRepos(filteredRepos.map((r) => r.name));
    }
  };

  const toggleSelectRepo = (name: string) => {
    setSelectedRepos((prev) =>
      prev.includes(name) ? prev.filter((r) => r !== name) : [...prev, name]
    );
  };

  const handleTransfer = async () => {
    if (!targetUser.trim()) {
      alert("Please specify a target account or organization.");
      return;
    }
    if (selectedRepos.length === 0) return;

    const confirmed = confirm(
      `Confirm transfer of ${selectedRepos.length} repository(ies) to '${targetUser}'? This will permanently transfer administrative ownership.`
    );
    if (!confirmed) return;

    setTransferring(true);
    addLog(`Initiating bulk transfer of ${selectedRepos.length} repo(s) to ${targetUser}...`, "info");

    const headers: Record<string, string> = { "Content-Type": "application/json" };
    if (patToken.trim()) {
      headers["Authorization"] = `Bearer ${patToken.trim()}`;
    }

    try {
      const res = await fetch("/api/transfer", {
        method: "POST",
        headers,
        body: JSON.stringify({
          target_user: targetUser.trim(),
          repositories: selectedRepos,
        }),
      });

      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Transfer failed");

      const transferResults: TransferResult[] = Array.isArray(data) ? data : data.results || [];
      setResults(transferResults);
      transferResults.forEach((item: TransferResult) => {
        if (item.success) {
          addLog(`Transferred ${item.repo} -> ${targetUser} (HTTP 202)`, "success");
        } else {
          addLog(`Failed ${item.repo}: ${item.message} (${item.status_code})`, "error");
        }
      });
      addLog("Batch execution completed.", "info");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Error executing transfer";
      addLog(`Fatal: ${msg}`, "error");
    } finally {
      setTransferring(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col bg-zinc-950">
      {/* Header */}
      <header className="border-b border-zinc-800 bg-zinc-900/60 backdrop-blur sticky top-0 z-50">
        <div className="max-w-6xl mx-auto px-6 h-16 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-lg overflow-hidden border border-blue-500/30 flex-shrink-0 bg-blue-950">
              <img src="/oauth-logo.png" alt="GHRA Logo" className="w-full h-full object-cover" />
            </div>
            <div>
              <span className="font-semibold text-zinc-100 text-sm tracking-tight">GitHub Repository Transfer</span>
              <span className="ml-2 text-xs font-mono bg-zinc-800 text-zinc-400 px-2 py-0.5 rounded border border-zinc-700">Next.js 16</span>
            </div>
          </div>

          <div>
            {user ? (
              <div className="flex items-center gap-3 bg-zinc-800/80 border border-zinc-700 px-3 py-1.5 rounded-full text-xs">
                <img src={user.avatar_url} alt={user.login} className="w-5 h-5 rounded-full" />
                <span className="font-medium text-zinc-200">{user.login}</span>
                <a href="/api/auth/logout" className="text-zinc-400 hover:text-red-400 ml-1 transition">
                  <LogOut className="w-3.5 h-3.5" />
                </a>
              </div>
            ) : (
              <a
                href="/api/auth/login"
                className="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg bg-blue-600 hover:bg-blue-500 text-white text-xs font-medium transition"
              >
                <GitPullRequest className="w-4 h-4" />
                Sign in with GitHub
              </a>
            )}
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-6xl w-full mx-auto px-6 py-8 flex-1 flex flex-col gap-6">
        {/* Banner if not logged in */}
        {!user && (
          <div className="bg-gradient-to-r from-blue-950/40 via-zinc-900 to-zinc-900 border border-blue-900/40 rounded-xl p-6 flex flex-col sm:flex-row items-center justify-between gap-4">
            <div>
              <h2 className="text-base font-semibold text-zinc-100 mb-1">One-Click GitHub Authentication</h2>
              <p className="text-xs text-zinc-400 max-w-xl leading-relaxed">
                Authenticate via GitHub OAuth to authorize repository transfers securely without generating or copying personal access tokens.
              </p>
            </div>
            <a
              href="/api/auth/login"
              className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white text-xs font-semibold rounded-lg shrink-0 transition flex items-center gap-2"
            >
              <GitPullRequest className="w-4 h-4" />
              Sign In with GitHub
            </a>
          </div>
        )}

        {/* Configuration Bar */}
        <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-5 grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-xs font-semibold text-zinc-400 uppercase tracking-wider mb-2">
              Destination GitHub Account / Org
            </label>
            <input
              type="text"
              value={targetUser}
              onChange={(e) => setTargetUser(e.target.value)}
              className="w-full bg-zinc-950 border border-zinc-800 focus:border-blue-500 text-sm text-zinc-200 px-3.5 py-2 rounded-lg outline-none font-mono"
              placeholder="e.g. MishraShardendu22"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-zinc-400 uppercase tracking-wider mb-2">
              Personal Access Token Fallback (Optional)
            </label>
            <input
              type="password"
              value={patToken}
              onChange={(e) => setPatToken(e.target.value)}
              className="w-full bg-zinc-950 border border-zinc-800 focus:border-blue-500 text-sm text-zinc-200 px-3.5 py-2 rounded-lg outline-none font-mono"
              placeholder="ghp_... (overrides OAuth session)"
            />
          </div>
        </div>

        {/* Repository Table Card */}
        <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-5 flex flex-col gap-4">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <div>
              <h3 className="text-sm font-semibold text-zinc-200">Select Repositories</h3>
              <p className="text-xs text-zinc-400">
                {repos.length > 0 ? `${filteredRepos.length} available (${selectedRepos.length} selected)` : "No repositories loaded"}
              </p>
            </div>

            <div className="flex items-center gap-2">
              <div className="relative">
                <Search className="w-4 h-4 absolute left-3 top-2.5 text-zinc-500" />
                <input
                  type="text"
                  placeholder="Search repositories..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="bg-zinc-950 border border-zinc-800 text-xs text-zinc-200 pl-9 pr-3 py-2 rounded-lg outline-none focus:border-blue-500 w-56 font-mono"
                />
              </div>

              <button
                onClick={fetchRepos}
                disabled={loading}
                className="p-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-lg text-xs font-medium transition disabled:opacity-50"
              >
                <RefreshCw className={`w-4 h-4 ${loading ? "animate-spin" : ""}`} />
              </button>

              <button
                onClick={handleTransfer}
                disabled={selectedRepos.length === 0 || transferring}
                className="px-4 py-2 bg-blue-600 hover:bg-blue-500 disabled:bg-zinc-800 disabled:text-zinc-600 text-white rounded-lg text-xs font-semibold transition"
              >
                {transferring ? "Transferring..." : `Transfer (${selectedRepos.length})`}
              </button>
            </div>
          </div>

          {/* Table */}
          <div className="border border-zinc-800 rounded-lg overflow-hidden bg-zinc-950 max-h-96 overflow-y-auto">
            <table className="w-full text-left border-collapse text-xs">
              <thead className="bg-zinc-900 border-b border-zinc-800 sticky top-0 z-10 text-zinc-400 uppercase font-mono">
                <tr>
                  <th className="p-3 w-10">
                    <input
                      type="checkbox"
                      checked={filteredRepos.length > 0 && selectedRepos.length === filteredRepos.length}
                      onChange={toggleSelectAll}
                      className="rounded bg-zinc-800 border-zinc-700 text-blue-600 focus:ring-0"
                    />
                  </th>
                  <th className="p-3">Repository</th>
                  <th className="p-3">Visibility</th>
                  <th className="p-3">Last Updated</th>
                  <th className="p-3 text-right">Transfer Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-800/60 font-sans">
                {filteredRepos.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="p-8 text-center text-zinc-500">
                      {loading ? "Loading repositories from GitHub..." : "No repositories found. Sign in or click refresh."}
                    </td>
                  </tr>
                ) : (
                  filteredRepos.map((repo) => {
                    const result = results.find((r) => r.repo === repo.name);
                    return (
                      <tr key={repo.name} className="hover:bg-zinc-900/40 transition">
                        <td className="p-3">
                          <input
                            type="checkbox"
                            checked={selectedRepos.includes(repo.name)}
                            onChange={() => toggleSelectRepo(repo.name)}
                            className="rounded bg-zinc-800 border-zinc-700 text-blue-600 focus:ring-0"
                          />
                        </td>
                        <td className="p-3">
                          <div className="flex items-center gap-2 font-mono text-zinc-200">
                            <FolderGit2 className="w-3.5 h-3.5 text-zinc-400 shrink-0" />
                            <span>{repo.name}</span>
                          </div>
                        </td>
                        <td className="p-3">
                          <span
                            className={`px-2 py-0.5 rounded text-[10px] font-mono uppercase font-semibold ${
                              repo.private
                                ? "bg-amber-950/40 text-amber-400 border border-amber-800/40"
                                : "bg-blue-950/40 text-blue-400 border border-blue-800/40"
                            }`}
                          >
                            {repo.private ? "Private" : "Public"}
                          </span>
                        </td>
                        <td className="p-3 text-zinc-400 font-mono">
                          {repo.updated_at ? new Date(repo.updated_at).toLocaleDateString() : "-"}
                        </td>
                        <td className="p-3 text-right">
                          {result ? (
                            result.success ? (
                              <span className="inline-flex items-center gap-1 text-emerald-400 font-mono text-[11px]">
                                <CheckCircle2 className="w-3.5 h-3.5" /> Accepted (202)
                              </span>
                            ) : (
                              <span className="inline-flex items-center gap-1 text-red-400 font-mono text-[11px]">
                                <AlertCircle className="w-3.5 h-3.5" /> Failed ({result.status_code})
                              </span>
                            )
                          ) : (
                            <span className="text-zinc-600 font-mono text-[11px]">Idle</span>
                          )}
                        </td>
                      </tr>
                    );
                  })
                )}
              </tbody>
            </table>
          </div>
        </div>

        {/* Live Terminal / Telemetry Console */}
        <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-5 flex flex-col gap-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Terminal className="w-4 h-4 text-zinc-400" />
              <h3 className="text-xs font-semibold text-zinc-300 font-mono uppercase tracking-wider">
                Execution Telemetry
              </h3>
            </div>
            <button
              onClick={() => setLogs([])}
              className="text-[11px] text-zinc-400 hover:text-zinc-200 transition"
            >
              Clear Log
            </button>
          </div>

          <div className="bg-black/90 border border-zinc-800/80 rounded-lg p-4 font-mono text-xs h-48 overflow-y-auto flex flex-col gap-1.5">
            {logs.map((log, i) => (
              <div key={i} className="flex gap-2">
                <span className="text-zinc-600 shrink-0">[{log.time}]</span>
                <span
                  className={
                    log.type === "success"
                      ? "text-emerald-400"
                      : log.type === "error"
                      ? "text-red-400"
                      : "text-blue-400"
                  }
                >
                  {log.text}
                </span>
              </div>
            ))}
          </div>
        </div>
      </main>
    </div>
  );
}
