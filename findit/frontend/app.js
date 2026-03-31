const API = 'http://localhost:8080';
let currentUser = null;
let allItems = [];

// --- Auth ---

function showTab(tab) {
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    if (tab === 'login') {
        document.querySelector('.tab:first-child').classList.add('active');
        document.getElementById('login-form').classList.remove('hidden');
        document.getElementById('signup-form').classList.add('hidden');
    } else {
        document.querySelector('.tab:last-child').classList.add('active');
        document.getElementById('login-form').classList.add('hidden');
        document.getElementById('signup-form').classList.remove('hidden');
    }
    clearAuthMessage();
}

function showMessage(el, text, type) {
    el.textContent = text;
    el.className = 'message ' + type;
    setTimeout(() => { el.textContent = ''; el.className = 'message'; }, 4000);
}

function clearAuthMessage() {
    const el = document.getElementById('auth-message');
    el.textContent = '';
    el.className = 'message';
}

async function handleSignup(e) {
    e.preventDefault();
    const data = {
        id: document.getElementById('signup-id').value.trim(),
        email: document.getElementById('signup-email').value.trim(),
        password: document.getElementById('signup-password').value
    };
    const msgEl = document.getElementById('auth-message');

    try {
        const res = await fetch(API + '/signup', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        const json = await res.json();
        if (res.ok) {
            showMessage(msgEl, 'Account created! You can now log in.', 'success');
            showTab('login');
        } else {
            showMessage(msgEl, json.error || 'Signup failed', 'error');
        }
    } catch (err) {
        showMessage(msgEl, 'Cannot connect to server', 'error');
    }
}

async function handleLogin(e) {
    e.preventDefault();
    const data = {
        email: document.getElementById('login-email').value.trim(),
        password: document.getElementById('login-password').value
    };
    const msgEl = document.getElementById('auth-message');

    try {
        const res = await fetch(API + '/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        const json = await res.json();
        if (res.ok) {
            currentUser = { email: data.email, user_id: json.user_id };
            showApp();
        } else {
            showMessage(msgEl, json.error || 'Login failed', 'error');
        }
    } catch (err) {
        showMessage(msgEl, 'Cannot connect to server', 'error');
    }
}

function logout() {
    currentUser = null;
    document.getElementById('auth-section').classList.remove('hidden');
    document.getElementById('app-section').classList.add('hidden');
    document.getElementById('login-email').value = '';
    document.getElementById('login-password').value = '';
}

// --- App ---

function showApp() {
    document.getElementById('auth-section').classList.add('hidden');
    document.getElementById('app-section').classList.remove('hidden');
    document.getElementById('user-info').textContent = currentUser.email;
    document.getElementById('item-date').valueAsDate = new Date();
    loadItems();
}

async function handleCreateItem(e) {
    e.preventDefault();
    const data = {
        id: document.getElementById('item-id').value.trim(),
        user_id: currentUser.user_id,
        type: document.getElementById('item-type').value,
        name: document.getElementById('item-name').value.trim(),
        description: document.getElementById('item-description').value.trim(),
        location: document.getElementById('item-location').value.trim(),
        date: document.getElementById('item-date').value
    };
    const msgEl = document.getElementById('item-message');

    try {
        const res = await fetch(API + '/items', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        const json = await res.json();
        if (res.ok) {
            showMessage(msgEl, 'Item reported successfully!', 'success');
            document.getElementById('item-form').reset();
            document.getElementById('item-date').valueAsDate = new Date();
            loadItems();
        } else {
            showMessage(msgEl, json.error || 'Failed to create item', 'error');
        }
    } catch (err) {
        showMessage(msgEl, 'Cannot connect to server', 'error');
    }
}

async function loadItems() {
    try {
        const res = await fetch(API + '/items');
        allItems = await res.json();
        renderItems(allItems);
    } catch (err) {
        document.getElementById('items-list').innerHTML = '<p class="empty">Cannot load items.</p>';
    }
}

function renderItems(items) {
    const list = document.getElementById('items-list');
    if (!items || items.length === 0) {
        list.innerHTML = '<p class="empty">No items reported yet.</p>';
        return;
    }
    list.innerHTML = items.map(item => `
        <div class="item-card ${item.type}">
            <div class="item-info">
                <h3>${escapeHtml(item.name)}</h3>
                <p>${escapeHtml(item.description || 'No description')}</p>
                <div class="item-meta">
                    📍 ${escapeHtml(item.location || 'Unknown')} &nbsp;|&nbsp; 📅 ${item.date || 'No date'} &nbsp;|&nbsp; 👤 ${escapeHtml(item.user_id)}
                </div>
            </div>
            <span class="item-badge ${item.type}">${item.type}</span>
        </div>
    `).join('');
}

function filterItems(type) {
    document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
    event.target.classList.add('active');
    if (type === 'all') {
        renderItems(allItems);
    } else {
        renderItems(allItems.filter(i => i.type === type));
    }
}

function escapeHtml(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

// --- Init ---
document.addEventListener('DOMContentLoaded', () => {
    // Set today's date as default
    document.getElementById('item-date').valueAsDate = new Date();
});
