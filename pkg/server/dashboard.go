package server

const DashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>GitHub Repository Transfer Engine</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet">
  <style>
    :root {
      color-scheme: dark;
      --bg-base: #09090b;
      --bg-surface: #121318;
      --bg-card: #181920;
      --bg-card-hover: #20222c;
      --border-default: #27272a;
      --border-focus: #3b82f6;
      --text-primary: #f4f4f5;
      --text-secondary: #a1a1aa;
      --text-muted: #71717a;
      --primary: #2563eb;
      --primary-hover: #1d4ed8;
      --accent: #3b82f6;
      --success: #10b981;
      --warning: #f59e0b;
      --danger: #ef4444;
      --font-sans: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      --font-mono: 'JetBrains Mono', monospace;
      --radius: 8px;
    }

    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: var(--bg-base);
      color: var(--text-primary);
      font-family: var(--font-sans);
      font-size: 14px;
      line-height: 1.5;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
    }

    header {
      background: var(--bg-surface);
      border-bottom: 1px solid var(--border-default);
      position: sticky;
      top: 0;
      z-index: 50;
    }

    .header-container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 16px 24px;
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .logo {
      display: flex;
      align-items: center;
      gap: 12px;
      font-weight: 700;
      font-size: 16px;
      letter-spacing: -0.01em;
      color: var(--text-primary);
      text-decoration: none;
    }

    .logo-badge {
      background: #1e293b;
      color: #60a5fa;
      padding: 2px 8px;
      border-radius: 4px;
      font-size: 12px;
      font-family: var(--font-mono);
      font-weight: 600;
      border: 1px solid #3b82f633;
    }

    .user-pill {
      display: flex;
      align-items: center;
      gap: 10px;
      background: var(--bg-card);
      border: 1px solid var(--border-default);
      padding: 6px 14px;
      border-radius: 9999px;
      font-size: 13px;
    }

    .user-avatar {
      width: 24px;
      height: 24px;
      border-radius: 50%;
      background: #27272a;
    }

    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 8px;
      padding: 8px 16px;
      font-size: 13px;
      font-weight: 600;
      border-radius: var(--radius);
      border: 1px solid transparent;
      cursor: pointer;
      transition: all 0.15s ease;
      font-family: var(--font-sans);
      text-decoration: none;
    }

    .btn-primary {
      background: var(--primary);
      color: #ffffff;
    }
    .btn-primary:hover:not(:disabled) {
      background: var(--primary-hover);
    }
    .btn-secondary {
      background: var(--bg-card);
      color: var(--text-primary);
      border-color: var(--border-default);
    }
    .btn-secondary:hover:not(:disabled) {
      background: var(--bg-card-hover);
      border-color: #3f3f46;
    }
    .btn-danger {
      background: rgba(239, 68, 68, 0.15);
      color: var(--danger);
      border-color: rgba(239, 68, 68, 0.3);
    }
    .btn-danger:hover:not(:disabled) {
      background: var(--danger);
      color: #ffffff;
    }
    .btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    main {
      max-width: 1200px;
      width: 100%;
      margin: 0 auto;
      padding: 32px 24px;
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 24px;
    }

    .card {
      background: var(--bg-surface);
      border: 1px solid var(--border-default);
      border-radius: var(--radius);
      padding: 24px;
    }

    .card-title {
      font-size: 16px;
      font-weight: 600;
      margin-bottom: 6px;
      color: var(--text-primary);
    }

    .card-desc {
      font-size: 13px;
      color: var(--text-secondary);
      margin-bottom: 20px;
    }

    .form-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 16px;
    }

    @media (max-width: 768px) {
      .form-grid {
        grid-template-columns: 1fr;
      }
    }

    .form-group {
      display: flex;
      flex-direction: column;
      gap: 6px;
    }

    .form-label {
      font-size: 12px;
      font-weight: 600;
      color: var(--text-secondary);
      text-transform: uppercase;
      letter-spacing: 0.04em;
    }

    .form-input {
      background: var(--bg-card);
      border: 1px solid var(--border-default);
      color: var(--text-primary);
      padding: 10px 14px;
      border-radius: var(--radius);
      font-size: 14px;
      font-family: inherit;
      outline: none;
      transition: border-color 0.15s ease;
    }

    .form-input:focus {
      border-color: var(--border-focus);
    }

    .repo-toolbar {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
      margin-bottom: 16px;
      flex-wrap: wrap;
    }

    .search-box {
      flex: 1;
      min-width: 240px;
    }

    .repo-table-wrapper {
      max-height: 420px;
      overflow-y: auto;
      border: 1px solid var(--border-default);
      border-radius: var(--radius);
      background: var(--bg-base);
    }

    table {
      width: 100%;
      border-collapse: collapse;
      text-align: left;
    }

    th {
      background: var(--bg-surface);
      color: var(--text-secondary);
      font-size: 12px;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      padding: 12px 16px;
      border-bottom: 1px solid var(--border-default);
      position: sticky;
      top: 0;
      z-index: 10;
    }

    td {
      padding: 12px 16px;
      border-bottom: 1px solid #1f1f23;
      font-size: 13px;
    }

    tr:hover td {
      background: rgba(255, 255, 255, 0.02);
    }

    .repo-name {
      font-weight: 600;
      font-family: var(--font-mono);
      color: var(--text-primary);
    }

    .badge {
      display: inline-block;
      padding: 2px 8px;
      border-radius: 4px;
      font-size: 11px;
      font-weight: 600;
      text-transform: uppercase;
      font-family: var(--font-mono);
    }

    .badge-public {
      background: rgba(59, 130, 246, 0.12);
      color: #60a5fa;
      border: 1px solid rgba(59, 130, 246, 0.25);
    }

    .badge-private {
      background: rgba(245, 158, 11, 0.12);
      color: #fbbf24;
      border: 1px solid rgba(245, 158, 11, 0.25);
    }

    .status-badge {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 2px 8px;
      border-radius: 4px;
      font-size: 11px;
      font-weight: 600;
      font-family: var(--font-mono);
    }
    .status-pending { background: #27272a; color: var(--text-secondary); }
    .status-in-progress { background: rgba(59, 130, 246, 0.15); color: #60a5fa; }
    .status-success { background: rgba(16, 185, 129, 0.15); color: var(--success); }
    .status-failed { background: rgba(239, 68, 68, 0.15); color: var(--danger); }

    .console-log {
      background: #000000;
      border: 1px solid var(--border-default);
      border-radius: var(--radius);
      padding: 16px;
      font-family: var(--font-mono);
      font-size: 12px;
      height: 220px;
      overflow-y: auto;
      display: flex;
      flex-direction: column;
      gap: 6px;
    }

    .log-line {
      display: flex;
      gap: 12px;
    }

    .log-time {
      color: var(--text-muted);
    }

    .log-info { color: #60a5fa; }
    .log-success { color: #34d399; }
    .log-error { color: #f87171; }
    .log-warn { color: #fbbf24; }

    .auth-banner {
      background: linear-gradient(135deg, #1e1b4b, #0f172a);
      border: 1px solid #3730a3;
      border-radius: var(--radius);
      padding: 32px;
      display: flex;
      flex-direction: column;
      align-items: center;
      text-align: center;
      gap: 16px;
    }

    .auth-title {
      font-size: 20px;
      font-weight: 700;
      color: #e0e7ff;
    }

    .auth-desc {
      color: #c7d2fe;
      max-width: 600px;
      font-size: 14px;
    }

    .hidden { display: none !important; }
  </style>
</head>
<body>

  <header>
    <div class="header-container">
      <a href="/" class="logo">
        <span>GitHub Repository Transfer Engine</span>
        <span class="logo-badge">REST API v2</span>
      </a>

      <div id="authNav">
        <!-- populated via js -->
      </div>
    </div>
  </header>

  <main>
    <!-- Authentication Card (when not logged in) -->
    <div id="authSection" class="auth-banner">
      <div class="auth-title">Seamless GitHub Account Authentication</div>
      <div class="auth-desc">
        Sign in via GitHub OAuth to transfer repositories without creating or managing Personal Access Tokens manually. You can also enter a Personal Access Token directly below for ad-hoc operations.
      </div>
      <div style="display: flex; gap: 12px; flex-wrap: wrap; justify-content: center; margin-top: 8px;">
        <a href="/api/auth/login" class="btn btn-primary" id="oauthBtn">
          Sign In with GitHub OAuth
        </a>
      </div>
    </div>

    <!-- Manual Token / Target Configuration Form -->
    <div class="card">
      <div class="card-title">Transfer Target & Authentication Settings</div>
      <div class="card-desc">Specify the destination GitHub account/organization and your credentials if not authenticated via OAuth.</div>

      <div class="form-grid">
        <div class="form-group">
          <label class="form-label" for="targetUser">Target Account or Organization</label>
          <input type="text" id="targetUser" class="form-input" placeholder="e.g. MishraShardendu22" value="MishraShardendu22">
        </div>

        <div class="form-group">
          <label class="form-label" for="patToken">Manual GitHub Token (Optional fallback)</label>
          <input type="password" id="patToken" class="form-input" placeholder="ghp_... or gho_... (uses OAuth session if blank)">
        </div>
      </div>
    </div>

    <!-- Repository Selection Table -->
    <div class="card">
      <div class="repo-toolbar">
        <div>
          <div class="card-title" style="margin-bottom: 2px;">Select Repositories to Transfer</div>
          <div id="repoCountText" style="font-size: 13px; color: var(--text-secondary);">Loading repositories...</div>
        </div>

        <div style="display: flex; gap: 10px; align-items: center;">
          <input type="text" id="repoSearch" class="form-input search-box" placeholder="Filter repositories...">
          <button id="refreshBtn" class="btn btn-secondary">Refresh</button>
          <button id="startTransferBtn" class="btn btn-primary" disabled>Transfer Selected Repos</button>
        </div>
      </div>

      <div class="repo-table-wrapper">
        <table>
          <thead>
            <tr>
              <th style="width: 40px;"><input type="checkbox" id="selectAllCheckbox"></th>
              <th>Repository</th>
              <th>Visibility</th>
              <th>Updated</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody id="repoTableBody">
            <tr>
              <td colspan="5" style="text-align: center; padding: 32px; color: var(--text-muted);">
                Sign in or click Refresh to list your repositories.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Live Execution Telemetry -->
    <div class="card">
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
        <div class="card-title" style="margin-bottom: 0;">Execution Telemetry & Console Log</div>
        <button id="clearLogBtn" class="btn btn-secondary" style="padding: 4px 10px; font-size: 12px;">Clear Console</button>
      </div>

      <div id="consoleLog" class="console-log">
        <div class="log-line"><span class="log-time">[SYSTEM]</span><span class="log-info">Ready. Connect your account or enter a token to initiate repository transfers.</span></div>
      </div>
    </div>
  </main>

  <script>
    let allRepos = [];
    let currentUser = null;
    let transferInProgress = false;

    const authSection = document.getElementById('authSection');
    const authNav = document.getElementById('authNav');
    const repoTableBody = document.getElementById('repoTableBody');
    const repoCountText = document.getElementById('repoCountText');
    const repoSearch = document.getElementById('repoSearch');
    const selectAllCheckbox = document.getElementById('selectAllCheckbox');
    const startTransferBtn = document.getElementById('startTransferBtn');
    const refreshBtn = document.getElementById('refreshBtn');
    const consoleLog = document.getElementById('consoleLog');
    const targetUserInput = document.getElementById('targetUser');
    const patTokenInput = document.getElementById('patToken');

    function log(msg, type = 'info') {
      const line = document.createElement('div');
      line.className = 'log-line';
      const time = new Date().toLocaleTimeString();
      line.innerHTML = '<span class="log-time">[' + time + ']</span><span class="log-' + type + '">' + escapeHtml(msg) + '</span>';
      consoleLog.appendChild(line);
      consoleLog.scrollTop = consoleLog.scrollHeight;
    }

    function escapeHtml(text) {
      const div = document.createElement('div');
      div.innerText = text;
      return div.innerHTML;
    }

    async function checkAuth() {
      try {
        const res = await fetch('/api/auth/me');
        if (res.ok) {
          currentUser = await res.json();
          authSection.classList.add('hidden');
          authNav.innerHTML = 
            '<div class="user-pill">' +
              '<img class="user-avatar" src="' + (currentUser.avatar_url || '') + '" alt="' + currentUser.login + '">' +
              '<span style="font-weight: 600;">' + currentUser.login + '</span>' +
              '<a href="/api/auth/logout" class="btn btn-secondary" style="padding: 2px 8px; font-size: 11px; margin-left: 6px;">Sign Out</a>' +
            '</div>';
          log('Authenticated as ' + currentUser.login + ' via GitHub OAuth.', 'success');
          loadRepos();
        } else {
          currentUser = null;
          authSection.classList.remove('hidden');
          authNav.innerHTML = '<a href="/api/auth/login" class="btn btn-primary">Sign In with GitHub</a>';
        }
      } catch (err) {
        console.error(err);
      }
    }

    async function loadRepos() {
      repoCountText.innerText = 'Loading repositories...';
      const token = patTokenInput.value.trim();
      const headers = {};
      if (token) headers['Authorization'] = 'Bearer ' + token;

      try {
        const res = await fetch('/api/repos', { headers });
        if (!res.ok) {
          const data = await res.json();
          throw new Error(data.error || 'Failed to fetch repositories');
        }
        allRepos = await res.json();
        renderRepos();
        log('Loaded ' + allRepos.length + ' repositories successfully.', 'info');
      } catch (err) {
        repoTableBody.innerHTML = '<tr><td colspan="5" style="text-align:center; padding:24px; color:var(--danger);">' + escapeHtml(err.message) + '</td></tr>';
        repoCountText.innerText = '0 repositories';
        log('Error loading repositories: ' + err.message, 'error');
      }
    }

    function renderRepos() {
      const filter = repoSearch.value.toLowerCase().trim();
      const filtered = allRepos.filter(r => r.name.toLowerCase().includes(filter));
      repoCountText.innerText = filtered.length + ' repositories found (' + allRepos.length + ' total)';

      if (filtered.length === 0) {
        repoTableBody.innerHTML = '<tr><td colspan="5" style="text-align:center; padding:24px; color:var(--text-muted);">No matching repositories found.</td></tr>';
        updateSelectionState();
        return;
      }

      let html = '';
      filtered.forEach(repo => {
        const isPrivate = repo.private;
        const badgeClass = isPrivate ? 'badge-private' : 'badge-public';
        const badgeText = isPrivate ? 'Private' : 'Public';
        const updatedDate = repo.updated_at ? new Date(repo.updated_at).toLocaleDateString() : '-';

        html += '<tr id="repo-row-' + repo.name + '">' +
          '<td><input type="checkbox" class="repo-select-cb" data-repo="' + repo.name + '"></td>' +
          '<td><span class="repo-name">' + escapeHtml(repo.name) + '</span></td>' +
          '<td><span class="badge ' + badgeClass + '">' + badgeText + '</span></td>' +
          '<td style="color:var(--text-secondary);">' + updatedDate + '</td>' +
          '<td><span class="status-badge status-pending" id="status-' + repo.name + '">Idle</span></td>' +
        '</tr>';
      });

      repoTableBody.innerHTML = html;

      document.querySelectorAll('.repo-select-cb').forEach(cb => {
        cb.addEventListener('change', updateSelectionState);
      });

      updateSelectionState();
    }

    function updateSelectionState() {
      const cbs = Array.from(document.querySelectorAll('.repo-select-cb'));
      const selected = cbs.filter(cb => cb.checked);
      startTransferBtn.disabled = selected.length === 0 || transferInProgress;
      startTransferBtn.innerText = 'Transfer ' + selected.length + ' Repos';

      if (cbs.length > 0 && selected.length === cbs.length) {
        selectAllCheckbox.checked = true;
        selectAllCheckbox.indeterminate = false;
      } else if (selected.length > 0) {
        selectAllCheckbox.checked = false;
        selectAllCheckbox.indeterminate = true;
      } else {
        selectAllCheckbox.checked = false;
        selectAllCheckbox.indeterminate = false;
      }
    }

    selectAllCheckbox.addEventListener('change', () => {
      const checked = selectAllCheckbox.checked;
      document.querySelectorAll('.repo-select-cb').forEach(cb => {
        cb.checked = checked;
      });
      updateSelectionState();
    });

    repoSearch.addEventListener('input', renderRepos);
    refreshBtn.addEventListener('click', loadRepos);
    document.getElementById('clearLogBtn').addEventListener('click', () => {
      consoleLog.innerHTML = '';
    });

    startTransferBtn.addEventListener('click', async () => {
      const targetUser = targetUserInput.value.trim();
      if (!targetUser) {
        alert('Please specify a target account or organization name.');
        return;
      }

      const selectedCbs = Array.from(document.querySelectorAll('.repo-select-cb')).filter(cb => cb.checked);
      const selectedRepos = selectedCbs.map(cb => cb.getAttribute('data-repo'));

      if (selectedRepos.length === 0) return;

      const confirmed = confirm('Are you sure you want to transfer ' + selectedRepos.length + ' repository(ies) to ' + targetUser + '? This is an administrative operation on GitHub.');
      if (!confirmed) return;

      transferInProgress = true;
      startTransferBtn.disabled = true;
      log('Starting bulk transfer of ' + selectedRepos.length + ' repositories to ' + targetUser + '...', 'info');

      selectedRepos.forEach(name => {
        const badge = document.getElementById('status-' + name);
        if (badge) {
          badge.className = 'status-badge status-in-progress';
          badge.innerText = 'Transferring...';
        }
      });

      const token = patTokenInput.value.trim();
      const headers = { 'Content-Type': 'application/json' };
      if (token) headers['Authorization'] = 'Bearer ' + token;

      try {
        const res = await fetch('/api/transfer', {
          method: 'POST',
          headers,
          body: JSON.stringify({
            target_user: targetUser,
            repositories: selectedRepos
          })
        });

        const data = await res.json();
        if (!res.ok) throw new Error(data.error || 'Transfer failed');

        data.results.forEach(resItem => {
          const badge = document.getElementById('status-' + resItem.repo);
          if (resItem.success) {
            if (badge) {
              badge.className = 'status-badge status-success';
              badge.innerText = 'Transferred (202)';
            }
            log('Transferred ' + resItem.repo + ' -> ' + targetUser + ' successfully.', 'success');
          } else {
            if (badge) {
              badge.className = 'status-badge status-failed';
              badge.innerText = 'Failed (' + resItem.status_code + ')';
            }
            log('Failed ' + resItem.repo + ': ' + resItem.message, 'error');
          }
        });

        log('Transfer batch completed.', 'info');
      } catch (err) {
        log('Fatal transfer error: ' + err.message, 'error');
        alert('Error: ' + err.message);
      } finally {
        transferInProgress = false;
        updateSelectionState();
      }
    });

    checkAuth();
  </script>
</body>
</html>
`
