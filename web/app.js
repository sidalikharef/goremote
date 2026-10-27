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
    discovered: [],
    discovering: false
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
    if (state.discovering) return;
    const list = document.getElementById('discoveredList');
    const btn = document.getElementById('startDiscoverBtn');
    if (list) {
      list.setAttribute('aria-live', 'polite');
      list.textContent = 'Searching for TVs...';
    }
    state.discovering = true;
    if (btn) btn.disabled = true;
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 10000);
      const res = await fetch('/api/discover', {method:'POST', credentials:'include', signal: controller.signal});
      clearTimeout(timeout);
      if (res.status === 401 || res.status === 403) {
        if (list) list.textContent = "Couldn't search for TVs. Authentication required.";
        return;
      }
      if (!res.ok) {
        if (list) list.textContent = "Couldn't search for TVs. Try again.";
        return;
      }
      const data = await res.json().catch(() => ({devices:[]}));
      state.discovered = data.devices || [];
      if (state.discovered.length === 0) {
        if (list) list.innerHTML = '<div>No TVs found.<br><small>Make sure your TV is powered on and connected to the same network, then try again.</small><br><button id="retryDiscover">Try again</button></div>';
        document.getElementById('retryDiscover')?.addEventListener('click', () => {
          if (list) list.textContent = '';
          discover();
        });
      } else {
        renderDiscovered();
      }
    } catch(e) {
      if (list) list.textContent = e.name === 'AbortError' ? "Couldn't search for TVs. Timed out. Try again." : "Couldn't search for TVs. Try again.";
    } finally {
      state.discovering = false;
      if (btn) btn.disabled = false;
    }
  }

  function renderDiscovered() {
    const list = document.getElementById('discoveredList');
    if (!list) return;
    list.innerHTML = '<div>Scanning for supported discoverable devices (Roku SSDP). Hisense VIDAA auto-discovery not verified.</div>';
    for (const d of state.discovered) {
      const div = document.createElement('div');
      div.className = 'disc-item';
      div.innerHTML = '<strong>' + (d.Name || d.name || d.ID || d.id) + '</strong> ' + (d.Address || d.address) + ' (' + (d.Driver || d.driver || 'roku') + ') <button data-id="' + (d.ID || d.id) + '" data-addr="' + (d.Address || d.address) + '" data-name="' + (d.Name || d.name) + '" data-driver="roku">Add</button>';
      list.appendChild(div);
    }
    list.querySelectorAll('button').forEach(btn => {
      btn.onclick = () => saveDiscovered(btn.getAttribute('data-id'), btn.getAttribute('data-addr'), btn.getAttribute('data-name'), btn.getAttribute('data-driver') || 'roku');
    });
  }

  async function saveDiscovered(id, addr, name, driver) {
    try {
      const res = await fetch('/api/discover/save', {method:'POST', headers:{'Content-Type':'application/json'}, credentials:'include', body: JSON.stringify({id:id, address:addr, name:name, driver:driver||'roku'})});
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

  function updateAuthUI() {
    const login = document.getElementById('authLogin');
    const enroll = document.getElementById('authEnroll');
    const logged = document.getElementById('authLogged');
    if (state.authed) {
      if (login) login.style.display = 'none';
      if (enroll) enroll.style.display = 'none';
      if (logged) logged.style.display = 'block';
    } else if (state.hasUsers) {
      if (login) login.style.display = 'block';
      if (enroll) enroll.style.display = 'none';
      if (logged) logged.style.display = 'none';
    } else {
      if (login) login.style.display = 'none';
      if (enroll) enroll.style.display = 'block';
      if (logged) logged.style.display = 'none';
    }
  }

  function updateAppsUI() {
    const list = document.getElementById('appsList');
    if (!list) return;
    list.innerHTML = '';
    for (const a of state.apps) {
      const btn = document.createElement('button');
      btn.textContent = a.title || a.name || a.id;
      btn.onclick = () => launchApp(a.id);
      list.appendChild(btn);
    }
  }

  function bindButtons() {
    document.getElementById('loginBtn')?.addEventListener('click', login);
    document.getElementById('logoutBtn')?.addEventListener('click', logout);
    document.getElementById('enrollBtn')?.addEventListener('click', enroll);
    document.getElementById('startDiscoverBtn')?.addEventListener('click', discover);
    document.getElementById('btnUp')?.addEventListener('click', () => sendKey('Up'));
    document.getElementById('btnDown')?.addEventListener('click', () => sendKey('Down'));
    document.getElementById('btnLeft')?.addEventListener('click', () => sendKey('Left'));
    document.getElementById('btnRight')?.addEventListener('click', () => sendKey('Right'));
    document.getElementById('btnSelect')?.addEventListener('click', () => sendKey('Select'));
    document.getElementById('btnBack')?.addEventListener('click', () => sendKey('Back'));
    document.getElementById('btnHome')?.addEventListener('click', () => sendKey('Home'));
    document.getElementById('btnPlay')?.addEventListener('click', () => sendKey('Play'));
    document.getElementById('btnPause')?.addEventListener('click', () => sendKey('Pause'));
    document.getElementById('btnVolUp')?.addEventListener('click', () => sendKey('VolumeUp'));
    document.getElementById('btnVolDown')?.addEventListener('click', () => sendKey('VolumeDown'));
    document.getElementById('btnMute')?.addEventListener('click', () => sendKey('VolumeMute'));
    document.getElementById('btnPower')?.addEventListener('click', () => sendKey('Power'));
  }

  async function vidaaRegister() {
    const id = document.getElementById('vidaaId')?.value || '';
    const name = document.getElementById('vidaaName')?.value || 'VIDAA TV';
    const addr = document.getElementById('vidaaAddr')?.value || '';
    const enable = document.getElementById('vidaaEnable')?.checked || false;
    try {
      const res = await fetch('/api/devices/vidaa/register', {method:'POST', headers:{'Content-Type':'application/json'}, credentials:'include', body: JSON.stringify({id:id, name:name, address:addr, enable:enable})});
      if (res.ok) {
        alert('Registered (experimental)');
        fetchDevices();
      } else {
        alert('Registration failed');
      }
    } catch(e) { alert('Error'); }
  }

  async function vidaaPairStart() {
    try {
      const res = await fetch('/api/devices/vidaa/pair/start', {method:'POST', headers:{'Content-Type':'application/json'}, credentials:'include', body: '{}'});
      if (res.ok) {
        const el = document.getElementById('vidaaPairMsg');
        if (el) el.textContent = 'Pairing started';
      } else {
        alert('Pair start failed');
      }
    } catch(e) { alert('Error'); }
  }

  async function vidaaPairSubmit() {
    const pin = document.getElementById('vidaaPin')?.value || '';
    if (pin.length !== 4) { alert('Enter 4 digits'); return; }
    try {
      const res = await fetch('/api/devices/vidaa/pair/submit', {method:'POST', headers:{'Content-Type':'application/json'}, credentials:'include', body: JSON.stringify({pin:pin})});
      if (res.ok) {
        const el = document.getElementById('vidaaPairMsg');
        if (el) el.textContent = 'Submitted';
      } else {
        alert('Pair failed');
      }
    } catch(e) { alert('Error'); }
  }

  async function vidaaPairCancel() {
    try {
      await fetch('/api/devices/vidaa/pair/cancel', {method:'POST', headers:{'Content-Type':'application/json'}, credentials:'include', body: '{}'});
    } catch(e) {}
  }

  async function vidaaUnpair() {
    try {
      await fetch('/api/devices/vidaa/unpair', {method:'POST', headers:{'Content-Type':'application/json'}, credentials:'include', body: '{}'});
      alert('Unpaired');
      fetchDevices();
    } catch(e) { alert('Error'); }
  }

  function bindVidaa() {
    document.getElementById('vidaaRegisterBtn')?.addEventListener('click', vidaaRegister);
    document.getElementById('vidaaPairStartBtn')?.addEventListener('click', vidaaPairStart);
    document.getElementById('vidaaPairSubmitBtn')?.addEventListener('click', vidaaPairSubmit);
    document.getElementById('vidaaPairCancelBtn')?.addEventListener('click', vidaaPairCancel);
    document.getElementById('vidaaUnpairBtn')?.addEventListener('click', vidaaUnpair);
  }

  function init() {
    checkAuth();
    bindButtons();
    bindVidaa();
    updateAuthUI();
    fetchDevices();
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
