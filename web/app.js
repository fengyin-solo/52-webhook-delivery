async function fetchJSON(url) {
  const res = await fetch(url);
  if (!res.ok) return null;
  return res.json();
}

function formatTime(iso) {
  if (!iso) return '-';
  const d = new Date(iso);
  return d.toLocaleString('zh-CN');
}

async function loadStats() {
  const data = await fetchJSON('/api/stats/overview');
  if (!data || !data.data) return;
  const s = data.data;
  const container = document.getElementById('stats');
  container.innerHTML = `
    <div class="stat-card"><div class="stat-value">${s.total_count}</div><div class="stat-label">总投递</div></div>
    <div class="stat-card"><div class="stat-value">${s.delivered_count}</div><div class="stat-label">已送达</div></div>
    <div class="stat-card"><div class="stat-value">${s.failed_count}</div><div class="stat-label">失败</div></div>
    <div class="stat-card"><div class="stat-value">${s.pending_count}</div><div class="stat-label">待处理</div></div>
    <div class="stat-card"><div class="stat-value">${s.retrying_count}</div><div class="stat-label">重试中</div></div>
    <div class="stat-card"><div class="stat-value">${s.success_rate}%</div><div class="stat-label">成功率</div></div>
  `;
}

async function loadEndpoints() {
  const data = await fetchJSON('/api/endpoints');
  if (!data || !data.data || !data.data.items) return;
  const tbody = document.querySelector('#endpoints-table tbody');
  tbody.innerHTML = '';
  data.data.items.forEach(e => {
    const tr = document.createElement('tr');
    tr.innerHTML = `
      <td>${e.id}</td>
      <td>${e.url}</td>
      <td><span class="badge badge-${e.status}">${e.status}</span></td>
      <td>${e.max_retries}</td>
      <td>${e.timeout_ms}</td>
      <td>${e.description || '-'}</td>
      <td>${formatTime(e.created_at)}</td>
    `;
    tbody.appendChild(tr);
  });
}

async function loadDeliveries() {
  const data = await fetchJSON('/api/deliveries');
  if (!data || !data.data || !data.data.items) return;
  const tbody = document.querySelector('#deliveries-table tbody');
  tbody.innerHTML = '';
  data.data.items.forEach(d => {
    const tr = document.createElement('tr');
    tr.innerHTML = `
      <td>${d.id}</td>
      <td>${d.event_id}</td>
      <td>${d.endpoint_id}</td>
      <td><span class="badge badge-${d.status}">${d.status}</span></td>
      <td>${d.attempts}</td>
      <td>${d.last_response_code || '-'}</td>
      <td>${formatTime(d.created_at)}</td>
    `;
    tbody.appendChild(tr);
  });
}

async function load() {
  await loadStats();
  await loadEndpoints();
  await loadDeliveries();
}

load();
