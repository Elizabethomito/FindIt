const API = '';
let currentUser = null;
let allItems = [];
let pollingInterval = null;
let currentItemId = null;

// --- Auth ---

function showTab(tab) {
    console.log('showTab called with tab:', tab);
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
    console.log('showMessage called with text:', text, 'type:', type);
    el.textContent = text;
    el.className = 'message ' + type;
    setTimeout(() => { el.textContent = ''; el.className = 'message'; }, 4000);
}

function clearAuthMessage() {
    console.log('clearAuthMessage called');
    const el = document.getElementById('auth-message');
    el.textContent = '';
    el.className = 'message';
}

async function handleSignup(e) {
    e.preventDefault();
    console.log('handleSignup called');
    const data = {
        id: document.getElementById('signup-id').value.trim(),
        email: document.getElementById('signup-email').value.trim(),
        password: document.getElementById('signup-password').value
    };
    console.log('Signup attempt for user:', data.id);
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
    console.log('handleLogin called');
    const data = {
        email: document.getElementById('login-email').value.trim(),
        password: document.getElementById('login-password').value
    };
    console.log('Login attempt for:', data.email);
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
    console.log('logout called');
    currentUser = null;
    document.getElementById('auth-section').classList.remove('hidden');
    document.getElementById('app-section').classList.add('hidden');
    document.getElementById('login-email').value = '';
    document.getElementById('login-password').value = '';
    stopPolling();
}

// --- App ---

function showApp() {
    console.log('showApp called');
    document.getElementById('auth-section').classList.add('hidden');
    document.getElementById('app-section').classList.remove('hidden');
    document.getElementById('user-info').textContent = currentUser.email;
    document.getElementById('item-date').valueAsDate = new Date();
    loadItems();
    startPolling();
}

async function handleCreateItem(e) {
    e.preventDefault();
    console.log('handleCreateItem called');
    const data = {
        user_id: currentUser.user_id,
        type: document.getElementById('item-type').value,
        name: document.getElementById('item-name').value.trim(),
        description: document.getElementById('item-description').value.trim(),
        location: document.getElementById('item-location').value.trim(),
        date: document.getElementById('item-date').value
    };
    console.log('Item data:', data);
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
    console.log('loadItems called');
    try {
        const res = await fetch(API + '/items');
        const newItems = await res.json();
        console.log('Loaded items:', newItems);
        
        // Check if items have changed
        if (JSON.stringify(newItems) !== JSON.stringify(allItems)) {
            allItems = newItems;
            renderItems(allItems);
        }
    } catch (err) {
        console.log('Error loading items:', err);
        document.getElementById('items-list').innerHTML = '<p class="empty">Cannot load items.</p>';
    }
}

function startPolling() {
    console.log('startPolling called');
    // Poll for new items every 3 seconds
    pollingInterval = setInterval(loadItems, 3000);
}

function stopPolling() {
    console.log('stopPolling called');
    if (pollingInterval) {
        clearInterval(pollingInterval);
        pollingInterval = null;
    }
}

function renderItems(items) {
    const list = document.getElementById('items-list');
    if (!items || items.length === 0) {
        list.innerHTML = '<p class="empty">No items reported yet.</p>';
        return;
    }
    console.log('Rendering items:', items);
    list.innerHTML = items.map(item => `
        <div class="item-card ${item.type}" onclick="showItemDetails('${item.id}')">
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
    console.log('filterItems called with type:', type);
    document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
    event.target.classList.add('active');
    if (type === 'all') {
        renderItems(allItems);
    } else {
        renderItems(allItems.filter(i => i.type === type));
    }
}

function escapeHtml(str) {
    console.log('escapeHtml called with str:', str);
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

// --- Item Details Modal ---

function showItemDetails(itemId) {
    console.log('showItemDetails called with itemId:', itemId);
    const item = allItems.find(i => i.id === itemId);
    console.log('Found item:', item);
    if (!item) {
        console.log('Item not found!');
        return;
    }
    
    currentItemId = itemId;
    document.getElementById('modal-name').textContent = item.name;
    document.getElementById('modal-type').textContent = item.type.charAt(0).toUpperCase() + item.type.slice(1);
    document.getElementById('modal-description').textContent = item.description || 'No description';
    document.getElementById('modal-location').textContent = item.location || 'Unknown';
    document.getElementById('modal-date').textContent = item.date || 'No date';
    document.getElementById('modal-user').textContent = item.user_id;
    
    // Show delete button only if current user is the item owner
    const deleteBtn = document.getElementById('delete-btn');
    if (currentUser && item.user_id === currentUser.user_id) {
        deleteBtn.classList.remove('hidden');
    } else {
        deleteBtn.classList.add('hidden');
    }
    
    const modal = document.getElementById('item-modal');
    console.log('Modal element:', modal);
    modal.classList.remove('hidden');
    console.log('Modal hidden class removed');
}

function closeModal() {
    console.log('closeModal called');
    document.getElementById('item-modal').classList.add('hidden');
    currentItemId = null;
}

async function deleteItem() {
    console.log('deleteItem called, currentItemId:', currentItemId);
    if (!currentItemId) return;
    
    if (!confirm('Are you sure you want to delete this item?')) {
        return;
    }
    
    try {
        const res = await fetch(API + '/items/' + currentItemId, {
            method: 'DELETE'
        });
        
        if (res.ok) {
            closeModal();
            loadItems(); // Reload items to reflect deletion
        } else {
            const json = await res.json();
            alert(json.error || 'Failed to delete item');
        }
    } catch (err) {
        alert('Cannot connect to server');
    }
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('item-modal');
    console.log('Window click event, target:', event.target, 'modal:', modal);
    if (event.target === modal) {
        closeModal();
    }
}

// --- Init ---
document.addEventListener('DOMContentLoaded', () => {
    console.log('DOMContentLoaded event fired');
    // Set today's date as default
    document.getElementById('item-date').valueAsDate = new Date();
});
