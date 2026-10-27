// GoRemote PWA frontend - vanilla JS
(function() {
  const state = {
    device: null,
    devices: [],
    apps: [],
    caps: {navigation:true, volume:true, power:true, media:true, apps:true, text_input:false, inputs:false, pairing:false},
    status: 'disconnected',
    authed: false,
    hasUsers: false,
    discovered: []
  };

  async function checkAuth() {
    try {
      const res = await fetch('/api/auth/status', {credentials:'include'});
      if (res.ok) {
        const data = await res.json();
        state.hasUsers = data.has_users;
        updateAuthUI();
        if (state.authed || state.hasUsers) {
          fetchDevices();
        }
      }
    } catch(e) {}
  }

  async function login() {
    const user = document.getElementById('loginUser')?.value || '';
    const passEl = document.getElementById('loginPass');
    const pass = passEl?.value || '';
    try {
      const res = await fetch('/api/auth/login', {method:'POST', headers:{'Content-Type':'application/json'}, credentials:'include', body: JSON.stringify({username:user, password:pass})});
      if (res.ok) {
        state.authed = true;
        if (passEl) passEl.value = '';
        updateAuthUI();
        fetchDevices();
      } else {
        alert('Login failed');
      }
    } catch(e) { alert('Login error'); }
  }

  async function logout() {
    try {
      await fetch('/api/auth/logout', {method:'POST', credentials:'include'});
    } catch(e) {}
    state.authed = false;
    updateAuthUI();
  }

  async function enroll() {
    const user = document.getElementById('enrollUser')?.value || '';
    const passEl = document.getElementById('enrollPass');
    const pass = passEl?.value || '';
    const token = document.getElementById('enrollToken')?.value || '';
    try {
      const res = await fetch('/api/auth/enroll', {method:'POST', headers:{'Content-Type':'application/json'}, credentials:'include', body: JSON.stringify({username:user, password:pass, bootstrap_token:token})});
      if (res.ok) {
        alert('Enrolled');
        if (passEl) passEl.value = '';
        checkAuth();
      } else {
        alert('Enrollment failed');
      }
    } catch(e) { alert('Enrollment error'); }
  }

  async function discover() {
    const list = document.getElementById('discoveredList');
    if (list) list.textContent = 'Discovering...';
    try {
      const res = await fetch('/api/discover', {method:'POST', credentials:'include'});
      if (res.ok) {
        const data = await res.json();
        state.discovered = data.devices || [];
        renderDiscovered();
      } else if (list) {
        list.textContent = 'Discovery failed or unauthorized';
      }
    } catch(e) {
      if (list) list.textContent = 'Discovery error';
    }
  }

  function renderDiscovered() {
    const list = document.getElementById('discoveredList');
    if (!list) return;
    list.innerHTML = '';
    for (const d of state.discovered) {
      const div = document.createElement('div');
      div.className = 'disc-item';
      div.innerHTML = '<strong>' + (d.Name || d.name || d.ID || d.id) + '</strong> ' + (d.Address || d.address) + ' <button data-id="' + (d.ID || d.id) + '" data-addr="' + (d.Address || d.address) + '" data-name="' + (d.Name || d.name) + '">Add</button>';
      list.appendChild(div);
    }
    list.querySelectorAll('button').forEach(btn => {
      btn.onclick = () => saveDiscovered(btn.getAttribute('data-id'), btn.getAttribute('data-addr'), btn.getAttribute('data-name'));
    });
  }

  async function saveDiscovered(id, addr, name) {
    try {
      const res = await fetch('/api/discover/save', {method:'POST', headers:{'Content-Type':'application/json'}, credentials:'include', body: JSON.stringify({id:id, address:addr, name:name, driver:'roku'})});
      if (res.ok || res.status === 201) {
        alert('Saved');
        fetchDevices();
      } else {
        alert('Save failed');
      }
    } catch(e) { alert('Save error'); }
  }

  async function fetchDevices() {
    try {
      const res = await fetch('/api/devices', {credentials:'include'});
      if (res.ok) {
        state.devices = await res.json();
        updateDeviceUI();
      } else if (res.status === 401) {
        state.authed = false;
        updateAuthUI();
      }
    } catch(e) {}
  }

  async function fetchApps() {
    if (!state.caps.apps) return;
    try {
      const sel = document.getElementById('deviceSelect');
      const dev = sel?.value || '';
      const res = await fetch('/api/apps' + (dev ? '?device=' + encodeURIComponent(dev) : ''), {credentials:'include'});
      if (res.ok) {
        state.apps = await res.json();
        updateAppsUI();
      }
    } catch(e) {}
  }

  async function sendKey(key) {
    setStatus('sending');
    try {
      const sel = document.getElementById('deviceSelect');
      const dev = sel?.value || '';
      const url = '/api/keypress/' + encodeURIComponent(key) + (dev ? '?device=' + encodeURIComponent(dev) : '');
      const res = await fetch(url, {method: 'POST', credentials:'include'});
      if (res.status === 401) { state.authed = false; updateAuthUI(); return; }
      if (res.status === 400 || res.status === 404 || res.status === 503) { setStatus('unavailable'); return; }
      if (!res.ok) throw new Error('failed');
      setStatus('connected');
    } catch(e) {
      setStatus('failure');
    }
  }

  async function launchApp(id) {
    try {
      const sel = document.getElementById('deviceSelect');
      const dev = sel?.value || '';
      const url = '/api/launch/' + encodeURIComponent(id) + (dev ? '?device=' + encodeURIComponent(dev) : '');
      await fetch(url, {method: 'POST', credentials:'include'});
    } catch(e) {}
  }

  function setStatus(s) {
    state.status = s;
    const el = document.getElementById('status');
    if (el) el.textContent = s;
  }

  function updateDeviceUI() {
    const sel = document.getElementById('deviceSelect');
    if (!sel) return;
    sel.innerHTML = '';
    for (const d of state.devices) {
      const opt = document.createElement('option');
      opt.value = d.id;
      opt.textContent = d.name + ' (' + d.driver + ') - ' + d.state;
      sel.appendChild(opt);
    }
  }

  function updateAppsUI() {
    const list = document.getElementById('appsList');
    if (!list) return;
    list.innerHTML = '';
    for (const a of state.apps) {
      const btn = document.createElement('button');
      btn.textContent = a.name || a.id;
      btn.className = 'app-btn';
      btn.onclick = () => launchApp(a.id);
      list.appendChild(btn);
    }
  }

  function updateAuthUI() {
    const authSec = document.getElementById('authSection');
    const loginSec = document.getElementById('loginSection');
    const enrollSec = document.getElementById('enrollSection');
    const remoteSec = document.getElementById('remoteSection');
    const appsSec = document.getElementById('appsSection');
    if (authSec) authSec.style.display = 'block';
    if (!state.hasUsers && enrollSec) { enrollSec.style.display = 'block'; if (loginSec) loginSec.style.display = 'none'; }
    if (state.hasUsers && loginSec) { loginSec.style.display = 'block'; if (enrollSec) enrollSec.style.display = 'none'; }
    if (state.authed) {
      if (remoteSec) remoteSec.style.display = 'block';
      if (appsSec) appsSec.style.display = 'block';
      if (loginSec) loginSec.style.display = 'none';
      if (enrollSec) enrollSec.style.display = 'none';
    } else {
      if (remoteSec) remoteSec.style.display = 'none';
      if (appsSec) appsSec.style.display = 'none';
    }
  }

  function init() {
    document.querySelectorAll('[data-key]').forEach(btn => {
      btn.addEventListener('click', () => {
        const key = btn.getAttribute('data-key');
        if (key) sendKey(key);
      });
    });

    document.getElementById('refreshApps')?.addEventListener('click', fetchApps);
    document.getElementById('loginBtn')?.addEventListener('click', login);
    document.getElementById('logoutBtn')?.addEventListener('click', logout);
    document.getElementById('enrollBtn')?.addEventListener('click', enroll);
    document.getElementById('startDiscoverBtn')?.addEventListener('click', discover);
    checkAuth();
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();

  function saveSelectedDevice() {
    const sel = document.getElementById('deviceSelect');
    if (sel && sel.value) {
      try {
        localStorage.setItem('goremote_selected_device', sel.value);
      } catch(e) {}
    }
  }

  function restoreSelectedDevice() {
    try {
      const v = localStorage.getItem('goremote_selected_device');
      if (v) {
        const sel = document.getElementById('deviceSelect');
        if (sel) sel.value = v;
      }
    } catch(e) {}
  }

  function init() {
    document.querySelectorAll('[data-key]').forEach(btn => {
      btn.addEventListener('click', () => {
        const key = btn.getAttribute('data-key');
        if (key) sendKey(key);
      });
    });

    document.getElementById('refreshApps')?.addEventListener('click', fetchApps);
    document.getElementById('loginBtn')?.addEventListener('click', login);
    document.getElementById('logoutBtn')?.addEventListener('click', logout);
    document.getElementById('enrollBtn')?.addEventListener('click', enroll);
    document.getElementById('startDiscoverBtn')?.addEventListener('click', discover);
    const sel = document.getElementById('deviceSelect');
    if (sel) {
      sel.addEventListener('change', saveSelectedDevice);
    }
    checkAuth();
    setTimeout(restoreSelectedDevice, 500);
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
