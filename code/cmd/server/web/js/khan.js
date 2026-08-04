/* ═══════════ خان — Khan Chat Frontend (Vue 3) ═══════════ */
const { createApp } = Vue;

// ─── دوزبانه: فارسی / English ───
const I18N = {
  fa: {
    tagline: 'جای گپ‌های تیم',
    about1: 'خان، چت سازمانی سبک برای شبکه محلی شماست. یک فایل، بدون اینترنت، بدون سرور جدا.',
    about2: 'امن، سریع و کاملا خصوصی — داده‌ها روی سیستم خودتان می‌ماند.',
    rolesTitle: 'نقش‌ها',
    roleUser: 'کاربر',
    roleUserDesc: 'چت و گفتگو + تغییر رمز خودش',
    roleSup: 'سوپروایزر',
    roleSupDesc: 'مدیریت گروه‌ها و اعضا',
    roleAdm: 'ادمین',
    roleAdmDesc: 'مدیریت کاربران و تنظیمات',
    featFast: 'سریع',
    featSecure: 'امن',
    featOffline: 'آفلاین',
    // Setup
    setupTitle: '⚙️ نصب اولیه',
    setupHint: 'سیستم برای اولین بار اجرا می‌شود. مدیر اصلی را بسازید:',
    setupUserLabel: 'نام کاربری مدیر',
    setupUserPlaceholder: 'مثلاً modir',
    setupNameLabel: 'نام نمایشی',
    setupNamePlaceholder: 'مثلاً شرکت آفتاب',
    setupPwdLabel: 'رمز عبور (حداقل ۸ کاراکتر)',
    setupBtn: '🚀 ساخت مدیر',
    setupBusy: 'در حال ساخت...',
    // Login
    loginTitle: 'ورود به خان',
    userLabel: 'نام کاربری',
    userPlaceholder: 'نام کاربری',
    pwdLabel: 'رمز عبور',
    pwdPlaceholder: '••••••••',
    loginBtn: 'ورود',
    loginBusy: 'در حال ورود...',
    credit: 'ساخته شده توسط aDiB 🧔',
    themeDark: 'حالت تاریک',
    themeLight: 'حالت روشن',
  },
  en: {
    tagline: 'Where Your Team Chats',
    about1: 'Khan is a lightweight team chat for your local network. One file, no internet, no separate server.',
    about2: 'Secure, fast and fully private — your data stays on your own machine.',
    rolesTitle: 'Roles',
    roleUser: 'User',
    roleUserDesc: 'Chat + change own password',
    roleSup: 'Supervisor',
    roleSupDesc: 'Manage groups and members',
    roleAdm: 'Admin',
    roleAdmDesc: 'Manage users and settings',
    featFast: 'Fast',
    featSecure: 'Secure',
    featOffline: 'Offline',
    // Setup
    setupTitle: '⚙️ Initial Setup',
    setupHint: 'First run detected. Create the main admin:',
    setupUserLabel: 'Admin username',
    setupUserPlaceholder: 'e.g. modir',
    setupNameLabel: 'Display name',
    setupNamePlaceholder: 'e.g. Aftab Company',
    setupPwdLabel: 'Password (min 8 chars)',
    setupBtn: '🚀 Create Admin',
    setupBusy: 'Creating...',
    // Login
    loginTitle: 'Login to Khan',
    userLabel: 'Username',
    userPlaceholder: 'Username',
    pwdLabel: 'Password',
    pwdPlaceholder: '••••••••',
    loginBtn: 'Login',
    loginBusy: 'Logging in...',
    credit: 'Built by aDiB 🧔',
    themeDark: 'Dark mode',
    themeLight: 'Light mode',
  },
};

const app = createApp({
  data() {
    return {
      // Language
      uiLang: localStorage.getItem('khan_lang') || 'fa',
      // Theme
      uiTheme: localStorage.getItem('khan_theme') || 'dark',
      i18n: I18N,
      // Server & session
      serverInfo: null,
      session: { token: localStorage.getItem('khan_token') || '', user: null },
      setupMode: false,
      setupForm: { admin_username: '', admin_password: '', company_name: '' },
      setupBusy: false,
      setupError: '',
      loginForm: { username: '', password: '' },
      loginBusy: false,
      loginError: '',

      // Data
      users: [],
      rooms: [],
      messages: {},
      roomMembers: {},
      currentRoom: null,
      currentRoomId: null,
      hasMore: {},
      loadingMore: false,

      // UI state
      searchQuery: '',
      mobileSidebar: false,
      // Unread tracking (client-side)
      readMap: {},   // roomId -> last read message id
      unreadMap: {}, // roomId -> count
      notifEnabled: false,
      docTitle: document.title,
      showNewChat: false,
      showCreateGroup: false,
      newGroupName: '',
      newChatSearch: '',
      showRoomInfo: false,
      showAddMember: false,
      memberSearch: '',
      showProfile: false,
      showChangePwd: false,
      pwdForm: { current: '', new1: '', new2: '' },
      pwdError: '',
      // Admin panel
      showAdminPanel: false,
      adminTab: 'users',
      showCreateUser: false,
      newUserForm: { username: '', display_name: '', password: '', role: 'user' },
      userFormError: '',
      showResetPwd: false,
      resetPwdTarget: {},
      resetPwdForm: '',
      license: { state: 'free', max_users: 20, licensed_to: '', expiry: '', error: '' },
      networkForm: { address_type: 'ip', ip: '', dns: '', port: 1727 },
      draft: '',
      editingMsg: null,
      emojiPicker: false,
      msgMenu: null,
      typingUsers: {},
      ws: null,
      wsRetry: 0,
      reconnectTimer: null,
      emojiList: ['😀','😁','😂','🤣','😊','😍','😘','😎','🤔','🙄','😅','😉','👍','👎','👏','🙏','🔥','❤️','💔','💯','🎉','🎊','😢','😭','🤝','✅','❌','⭐','🚀','💪'],
    };
  },

  computed: {
    filteredChats() {
      if (!this.searchQuery.trim()) return this.rooms;
      const q = this.searchQuery.toLowerCase();
      return this.rooms.filter(r => this.chatName(r).toLowerCase().includes(q));
    },

    totalUnread() {
      let total = 0;
      for (const k in this.unreadMap) total += this.unreadMap[k] || 0;
      return total;
    },

    roomStatus() {
      if (!this.currentRoom) return '';
      if (this.isDM(this.currentRoom)) {
        const other = this.dmPartner(this.currentRoom);
        if (other) return this.isOnline(other) ? 'آنلاین' : 'آفلاین';
        return '';
      }
      const count = this.currentRoomMembers.length;
      return count ? `${count} عضو` : '';
    },

    currentRoomMembers() {
      if (!this.currentRoom) return [];
      return this.roomMembers[this.currentRoom.id] || [];
    },

    groupedMessages() {
      const roomId = this.currentRoomId;
      const msgs = this.messages[roomId] || [];
      const groups = {};
      for (const m of msgs) {
        const key = this.dayKey(m.created_at);
        if (!groups[key]) groups[key] = [];
        groups[key].push(m);
      }
      return groups;
    },

    typingText() {
      const names = Object.values(this.typingUsers).filter(n => n !== this.session.user.id);
      if (!names.length) return '';
      const list = Object.keys(this.typingUsers).filter(id => id != this.session.user.id);
      const label = list.map(id => {
        const u = this.users.find(x => x.id == id);
        return u ? (u.display_name || u.username) : 'کسی';
      });
      return label.join('، ') + ' در حال نوشتن...';
    },

    typingUsersList() {
      return Object.keys(this.typingUsers).filter(id => id != this.session.user.id);
    },

    canManageRoom() {
      const u = this.session.user;
      if (!u || !this.currentRoom) return false;
      if (['adm', 'sadm'].includes(u.role)) return true;
      if (u.role === 'sup') {
        // supervisor manages groups they're a member of
        return this.currentRoom.type !== 'dm';
      }
      return false;
    },

    isAdmin() {
      const r = this.session.user ? this.session.user.role : '';
      return ['adm', 'sadm'].includes(r);
    },

    licenseIcon() {
      return { free: '🆓', valid: '✅', tampered: '⚠️' }[this.license.state] || '🎫';
    },

    licenseTitle() {
      return { free: 'نسخه رایگان', valid: 'لایسنس معتبر', tampered: 'لایسنس نامعتبر!' }[this.license.state] || 'لایسنس';
    },

    allVisibleUsers() {
      const q = (this.newChatSearch || '').toLowerCase();
      const me = this.session.user ? this.session.user.id : -1;
      return this.users.filter(u => u.id !== me && (!q || (u.display_name || u.username).toLowerCase().includes(q)));
    },

    searchableUsers() {
      const q = (this.memberSearch || '').toLowerCase();
      if (!this.currentRoom) return [];
      const existing = (this.roomMembers[this.currentRoom.id] || []).map(m => m.id);
      return this.users.filter(u =>
        !existing.includes(u.id) &&
        (!q || (u.display_name || u.username).toLowerCase().includes(q))
      );
    },
  },

  methods: {
    /* ─────────── i18n helper ─────────── */
    t(key) {
      return (this.i18n[this.uiLang] && this.i18n[this.uiLang][key]) || this.i18n.fa[key] || key;
    },

    // Method (takes param): admin can promote/demote supervisors
    canPromote(u) {
      if (!u || !this.session.user) return false;
      return this.isAdmin && u.role !== 'adm' && u.role !== 'sadm' && u.id !== this.session.user.id;
    },
    setLang(lang) {
      this.uiLang = lang;
      localStorage.setItem('khan_lang', lang);
    },

    /* ─────────── Theme ─────────── */
    toggleTheme() {
      this.uiTheme = this.uiTheme === 'dark' ? 'light' : 'dark';
      localStorage.setItem('khan_theme', this.uiTheme);
      document.documentElement.setAttribute('data-theme', this.uiTheme);
    },

    /* ─────────── API helpers ─────────── */
    async api(path, options = {}) {
      const headers = { 'Content-Type': 'application/json' };
      if (this.session.token) headers['Authorization'] = 'Bearer ' + this.session.token;
      const res = await fetch(path, {
        ...options,
        headers: { ...headers, ...(options.headers || {}) },
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(data.error || 'خطای سرور');
      return data;
    },

    /* ─────────── Setup / Login ─────────── */
    async checkSetup() {
      try {
        const data = await this.api('/api/setup/needs-setup');
        this.setupMode = data.needs_setup;
        if (!this.setupMode && this.session.token) {
          await this.loadSession();
        }
      } catch (e) {
        this.setupMode = true;
      }
    },

    async loadSession() {
      try {
        const me = await this.api('/api/auth/me');
        this.session.user = me;
        if (me.must_change_pwd) {
          // Force password change
          this.showChangePwd = true;
          this.pwdError = 'برای ادامه باید رمز عبور خود را تغییر دهید';
        }
        await this.loadAll();
        this.connectWS();
      } catch (e) {
        // Token invalid — back to login
        localStorage.removeItem('khan_token');
        this.session.token = '';
        this.session.user = null;
      }
    },

    async doSetup() {
      this.setupBusy = true;
      this.setupError = '';
      try {
        const data = await this.api('/api/setup', {
          method: 'POST',
          body: JSON.stringify(this.setupForm),
        });
        if (data.ok) {
          // Auto-login as admin
          this.loginForm.username = this.setupForm.admin_username;
          this.loginForm.password = this.setupForm.admin_password;
          await this.doLogin();
        }
      } catch (e) {
        this.setupError = e.message;
      } finally {
        this.setupBusy = false;
      }
    },

    async doLogin() {
      this.loginBusy = true;
      this.loginError = '';
      try {
        const data = await this.api('/api/auth/login', {
          method: 'POST',
          body: JSON.stringify(this.loginForm),
        });
        this.session.token = data.token;
        this.session.user = data.user;
        localStorage.setItem('khan_token', data.token);
        this.loginForm.password = '';
        await this.loadAll();
        this.connectWS();
      } catch (e) {
        this.loginError = e.message;
      } finally {
        this.loginBusy = false;
      }
    },

    async logout() {
      try { await this.api('/api/auth/logout', { method: 'POST' }); } catch (e) {}
      if (this.ws) this.ws.close();
      localStorage.removeItem('khan_token');
      this.session = { token: '', user: null };
      this.rooms = [];
      this.messages = {};
      this.unreadMap = {};
      this.readMap = {};
      this.currentRoom = null;
      this.currentRoomId = null;
      document.title = this.docTitle;
    },

    /* ─────────── Data loading ─────────── */
    async loadAll() {
      try {
        const [users, rooms] = await Promise.all([
          this.api('/api/users'),
          this.api('/api/rooms'),
        ]);
        // Defensive: backend may return null for empty lists (JSON)
        this.users = Array.isArray(users) ? users : [];
        this.rooms = Array.isArray(rooms) ? rooms : [];
        // Restore read map
        try {
          this.readMap = JSON.parse(localStorage.getItem('khan_readmap') || '{}');
        } catch (e) { this.readMap = {}; }
        // Load members for each room
        for (const room of this.rooms) {
          this.loadMembers(room.id).catch(() => {});
          this.loadMessages(room.id).catch(() => {});
        }
      } catch (e) {
        console.error('loadAll error:', e);
      }
    },

    async loadMembers(roomId) {
      try {
        const members = await this.api(`/api/rooms/${roomId}/members`);
        this.roomMembers[roomId] = Array.isArray(members) ? members : [];
      } catch (e) {}
    },

    async loadMessages(roomId) {
      try {
        const msgs = await this.api(`/api/messages/${roomId}?limit=50`);
        this.messages[roomId] = Array.isArray(msgs) ? msgs : [];
        this.hasMore[roomId] = (msgs && msgs.length) ? msgs.length >= 50 : false;
        // Compute unread count client-side
        this.computeUnread(roomId);
      } catch (e) {}
    },

    // Unread logic: count messages newer than the last-read message id
    computeUnread(roomId) {
      if (roomId === this.currentRoomId) {
        this.unreadMap[roomId] = 0;
        return;
      }
      const msgs = this.messages[roomId] || [];
      const lastRead = this.readMap[roomId] || 0;
      let count = 0;
      for (const m of msgs) {
        if (m.id > lastRead) count++;
      }
      this.unreadMap[roomId] = count;
    },

    async loadMore() {
      const roomId = this.currentRoomId;
      if (!roomId || this.loadingMore || !this.hasMore[roomId]) return;
      this.loadingMore = true;
      try {
        const msgs = this.messages[roomId] || [];
        const before = msgs.length ? msgs[0].id : 0;
        const older = await this.api(`/api/messages/${roomId}?limit=50&before=${before}`);
        if (older.length) {
          this.messages[roomId] = [...older, ...msgs];
        }
        this.hasMore[roomId] = older.length >= 50;
      } catch (e) {}
      this.loadingMore = false;
    },

    /* ─────────── Rooms ─────────── */
    chatName(room) {
      if (room.type === 'dm') {
        const other = this.dmPartner(room);
        return other ? (other.display_name || other.username) : 'گفتگو';
      }
      return room.name || 'گروه';
    },

    dmPartner(room) {
      if (!room || room.type !== 'dm' || !room.id) return null;
      // DM rooms have no explicit partner — infer from members
      const members = this.roomMembers[room.id] || [];
      const me = this.session.user ? this.session.user.id : -1;
      return members.find(m => m.id !== me) || null;
    },

    isDM(room) { return room && room.type === 'dm'; },

    roomTypeLabel(type) {
      return { dm: 'گفتگوی خصوصی', group: 'گروه', public: 'عمومی', private: 'خصوصی' }[type] || type;
    },

    chatPreview(room) {
      const msgs = this.messages[room.id] || [];
      if (!msgs.length) return this.roomTypeLabel(room.type);
      const last = msgs[msgs.length - 1];
      const sender = last.sender_id === this.session.user.id ? 'شما: ' :
        (last.sender_name ? last.sender_name + ': ' : '');
      return sender + (last.text || '📎 فایل');
    },

    chatTime(dateStr) {
      if (!dateStr) return '';
      const d = new Date(dateStr);
      const now = new Date();
      if (d.toDateString() === now.toDateString()) {
        return d.toLocaleTimeString('fa-IR', { hour: '2-digit', minute: '2-digit' });
      }
      return d.toLocaleDateString('fa-IR', { month: 'short', day: 'numeric' });
    },

    msgTime(dateStr) {
      if (!dateStr) return '';
      return new Date(dateStr).toLocaleTimeString('fa-IR', { hour: '2-digit', minute: '2-digit' });
    },

    dayKey(dateStr) {
      if (!dateStr) return '';
      const d = new Date(dateStr);
      const now = new Date();
      const yesterday = new Date(now); yesterday.setDate(now.getDate() - 1);
      if (d.toDateString() === now.toDateString()) return 'امروز';
      if (d.toDateString() === yesterday.toDateString()) return 'دیروز';
      return d.toLocaleDateString('fa-IR', { weekday: 'long', day: 'numeric', month: 'long' });
    },

    async openRoom(room) {
      this.currentRoom = room;
      this.currentRoomId = room.id;
      this.mobileSidebar = false;
      // Mark as read
      const msgs = this.messages[room.id] || [];
      const lastId = msgs.length ? msgs[msgs.length - 1].id : 0;
      this.readMap[room.id] = lastId;
      this.unreadMap[room.id] = 0;
      try { localStorage.setItem('khan_readmap', JSON.stringify(this.readMap)); } catch (e) {}
      this.updateDocTitle();
      // Refresh messages
      await this.loadMessages(room.id);
      await this.loadMembers(room.id);
      this.$nextTick(() => this.scrollToBottom());
    },

    async createGroup() {
      try {
        const room = await this.api('/api/rooms', {
          method: 'POST',
          body: JSON.stringify({ name: this.newGroupName.trim(), type: 'public' }),
        });
        this.rooms.unshift(room);
        this.newGroupName = '';
        this.showCreateGroup = false;
        this.showNewChat = false;
        await this.openRoom(room);
      } catch (e) {
        alert(e.message);
      }
    },

    async startDM(user) {
      try {
        const data = await this.api(`/api/rooms/dm/${user.id}`, { method: 'POST' });
        const room = this.rooms.find(r => r.id === data.room_id);
        if (room) {
          this.showNewChat = false;
          await this.openRoom(room);
        } else {
          // Refresh rooms
          await this.loadAll();
          const fresh = this.rooms.find(r => r.id === data.room_id);
          if (fresh) { this.showNewChat = false; await this.openRoom(fresh); }
        }
      } catch (e) {
        alert(e.message);
      }
    },

    async addMember(user) {
      try {
        await this.api(`/api/rooms/${this.currentRoom.id}/members`, {
          method: 'POST',
          body: JSON.stringify({ user_id: user.id }),
        });
        await this.loadMembers(this.currentRoom.id);
        this.memberSearch = '';
      } catch (e) {
        alert(e.message);
      }
    },

    async removeMember(user) {
      if (!confirm(`آیا ${user.display_name || user.username} از گروه حذف شود؟`)) return;
      try {
        await this.api(`/api/rooms/${this.currentRoom.id}/members/${user.id}`, { method: 'DELETE' });
        await this.loadMembers(this.currentRoom.id);
      } catch (e) {
        alert(e.message);
      }
    },

    /* ─────────── Messages ─────────── */
    sendMessage() {
      const text = this.draft.trim();
      if (!text && !this.editingMsg) return;

      if (this.editingMsg) {
        this.saveEdit();
        return;
      }

      if (!this.currentRoom) return;
      this.wsSend({ type: 'send_message', room_id: this.currentRoom.id, text });
      this.draft = '';
      this.autoGrow();
    },

    onTyping() {
      if (!this.currentRoom) return;
      this.wsSend({ type: 'typing', room_id: this.currentRoom.id });
    },

    async editMessageFromWs() {},

    startEdit(msg) {
      this.editingMsg = msg;
      this.draft = msg.text;
      this.msgMenu = null;
      this.$nextTick(() => {
        const ta = document.querySelector('.composer-input');
        if (ta) ta.focus();
      });
    },

    cancelEdit() {
      this.editingMsg = null;
      this.draft = '';
    },

    saveEdit() {
      if (!this.editingMsg) return;
      this.wsSend({ type: 'edit_message', message_id: this.editingMsg.id, text: this.draft.trim() });
      this.editingMsg = null;
      this.draft = '';
    },

    async deleteMessage(msg) {
      if (!confirm('حذف شود؟')) return;
      this.msgMenu = null;
      try {
        await this.api(`/api/messages/${msg.id}`, { method: 'DELETE' });
        const roomId = this.currentRoomId;
        this.messages[roomId] = (this.messages[roomId] || []).filter(m => m.id !== msg.id);
      } catch (e) {
        alert(e.message);
      }
    },

    canDeleteMessage(msg) {
      if (msg.sender_id === this.session.user.id) return true;
      return this.canManageRoom;
    },

    /* Reactions */
    aggregateReactions(msg) {
      const agg = {};
      for (const r of (msg.reactions || [])) {
        agg[r.emoji] = (agg[r.emoji] || 0) + 1;
      }
      return agg;
    },

    reactedByMe(msg, emoji) {
      return (msg.reactions || []).some(r => r.emoji === emoji && r.user_id === this.session.user.id);
    },

    async toggleReaction(msg, emoji) {
      const reacted = this.reactedByMe(msg, emoji);
      if (reacted) {
        this.wsSend({ type: 'remove_reaction', message_id: msg.id, emoji });
      } else {
        this.wsSend({ type: 'add_reaction', message_id: msg.id, emoji });
      }
    },

    quickReact(msg, emoji) {
      this.wsSend({ type: 'add_reaction', message_id: msg.id, emoji });
      this.msgMenu = null;
    },

    /* ─────────── Context menu ─────────── */
    openMsgMenu(msg, event) {
      this.msgMenu = { msg, x: event.clientX, y: event.clientY };
    },

    msgMenuStyle() {
      if (!this.msgMenu) return {};
      const menuW = 180, menuH = 180;
      let x = this.msgMenu.x;
      let y = this.msgMenu.y;
      if (x + menuW > window.innerWidth) x = window.innerWidth - menuW - 8;
      if (y + menuH > window.innerHeight) y = window.innerHeight - menuH - 8;
      return { left: x + 'px', top: y + 'px', position: 'fixed' };
    },

    /* ─────────── Emoji ─────────── */
    addEmoji(e) {
      this.draft += e;
      this.autoGrow();
    },

    /* ─────────── WebSocket ─────────── */
    connectWS() {
      if (!this.session.token) return;
      if (this.ws) { try { this.ws.close(); } catch(e) {} }
      const proto = location.protocol === 'https:' ? 'wss' : 'ws';
      this.ws = new WebSocket(`${proto}://${location.host}/ws?token=${encodeURIComponent(this.session.token)}`);

      this.ws.onopen = () => {
        this.wsRetry = 0;
        console.log('✅ Khan WS connected');
      };

      this.ws.onmessage = (e) => {
        try {
          const ev = JSON.parse(e.data);
          this.handleEvent(ev);
        } catch (err) {
          console.error('bad ws message', err);
        }
      };

      this.ws.onclose = () => {
        console.log('ws closed');
        // Reconnect with backoff
        if (this.session.token) {
          clearTimeout(this.reconnectTimer);
          this.reconnectTimer = setTimeout(() => this.connectWS(), Math.min(30000, 1000 * Math.pow(2, this.wsRetry++)));
        }
      };

      this.ws.onerror = (e) => { console.error('ws error', e); };
    },

    wsSend(obj) {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify(obj));
      }
    },

    handleEvent(ev) {
      switch (ev.type) {
        case 'message': {
          const roomId = ev.room_id;
          if (!this.messages[roomId]) this.messages[roomId] = [];
          this.messages[roomId].push(ev.payload);
          if (this.currentRoomId === roomId) {
            this.scrollToBottom();
            this.readMap[roomId] = ev.payload.id;
            this.unreadMap[roomId] = 0;
          } else {
            this.unreadMap[roomId] = (this.unreadMap[roomId] || 0) + 1;
            this.notifyNewMessage(ev.payload, roomId);
          }
          // Ensure room exists in list (for new DMs/groups)
          if (!this.rooms.some(r => r.id === roomId)) {
            this.loadAll();
          } else {
            // Move room to top (recent activity)
            const room = this.rooms.find(r => r.id === roomId);
            if (room) {
              this.rooms = this.rooms.filter(r => r.id !== roomId);
              this.rooms.unshift(room);
            }
          }
          this.updateDocTitle();
          break;
        }
        case 'message_edited': {
          const roomId = ev.room_id;
          const msgs = this.messages[roomId] || [];
          const idx = msgs.findIndex(m => m.id === ev.payload.id);
          if (idx >= 0) {
            msgs[idx].text = ev.payload.text;
            msgs[idx].edited_at = ev.payload.edited_at;
          }
          break;
        }
        case 'message_deleted': {
          const roomId = ev.room_id;
          this.messages[roomId] = (this.messages[roomId] || []).filter(m => m.id !== ev.payload);
          break;
        }
        case 'reaction': {
          const roomId = ev.room_id;
          const msgs = this.messages[roomId] || [];
          const msg = msgs.find(m => m.id === ev.payload.message_id);
          if (!msg) break;
          if (!msg.reactions) msg.reactions = [];
          if (ev.payload.remove) {
            msg.reactions = msg.reactions.filter(r => !(r.user_id === ev.payload.user_id && r.emoji === ev.payload.emoji));
          } else {
            if (!msg.reactions.some(r => r.user_id === ev.payload.user_id && r.emoji === ev.payload.emoji)) {
              msg.reactions.push({ user_id: ev.payload.user_id, emoji: ev.payload.emoji });
            }
          }
          break;
        }
        case 'typing': {
          const roomId = ev.room_id;
          if (roomId !== this.currentRoomId) break;
          const uid = ev.payload.user_id;
          if (uid === this.session.user.id) break;
          this.typingUsers[uid] = true;
          clearTimeout(this.typingUsers[uid + '_t']);
          this.typingUsers[uid + '_t'] = setTimeout(() => delete this.typingUsers[uid], 2500);
          break;
        }
        case 'presence': {
          // update user presence (visual only)
          const p = ev.payload;
          const u = this.users.find(x => x.id === p.user_id);
          if (u) u.online = p.online;
          break;
        }
        case 'force_logout':
          alert('شما از سیستم خارج شدید');
          this.logout();
          break;
        case 'error':
          alert(ev.payload.error || 'خطا');
          break;
      }
    },

    /* ─────────── Notifications & Unread ─────────── */
    updateDocTitle() {
      const total = this.totalUnread;
      if (total > 0) {
        document.title = `(${total}) ${this.docTitle}`;
      } else {
        document.title = this.docTitle;
      }
    },

    requestNotifPermission() {
      if ('Notification' in window && Notification.permission === 'default') {
        Notification.requestPermission().then(p => {
          this.notifEnabled = p === 'granted';
        });
      } else if ('Notification' in window) {
        this.notifEnabled = Notification.permission === 'granted';
      }
    },

    notifyNewMessage(msg, roomId) {
      // Browser notification (only when tab hidden or permission granted)
      if ('Notification' in window && Notification.permission === 'granted' && document.hidden) {
        try {
          const room = this.rooms.find(r => r.id === roomId);
          const title = this.chatName(room || { type: 'group' });
          const body = (msg.sender_name ? msg.sender_name + ': ' : '') + (msg.text || '📎 فایل');
          const n = new Notification(title, { body, icon: '/img/khan-logo.jpg', tag: 'khan-' + roomId });
          n.onclick = () => {
            window.focus();
            if (room) this.openRoom(room);
            n.close();
          };
        } catch (e) {}
      }
    },

    /* ─────────── Helpers ─────────── */
    isOnline(user) {
      return !!user.online;
    },

    avatarStyle(user) {
      if (!user) return { background: '#4d9de0' };
      const colors = ['#e17055', '#0984e3', '#00b894', '#6c5ce7', '#fdcb6e', '#e84393', '#00cec9', '#d63031'];
      const name = (user.display_name || user.username || '?');
      let hash = 0;
      for (const ch of name) hash = (hash * 31 + ch.charCodeAt(0)) % 997;
      return { background: colors[hash % colors.length] };
    },

    chatAvatarStyle(room) {
      if (room.type === 'dm') {
        const other = this.dmPartner(room);
        return this.avatarStyle(other || {});
      }
      return { background: 'linear-gradient(135deg, #4d9de0, #2b5278)' };
    },

    initials(name) {
      if (!name) return '?';
      const parts = name.trim().split(/\s+/);
      if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase();
      return name[0].toUpperCase();
    },

    safeUser(user) {
      return user || {};
    },

    roleLabel(role) {
      return { user: 'کاربر', sup: 'سوپروایزر', adm: 'ادمین', sadm: 'مدیر اصلی' }[role] || role;
    },

    scrollToBottom() {
      this.$nextTick(() => {
        const area = this.$refs.messagesArea;
        if (area) area.scrollTop = area.scrollHeight;
      });
    },

    onScroll() {
      const area = this.$refs.messagesArea;
      if (area && area.scrollTop < 40) {
        this.loadMore();
      }
    },

    autoGrow() {
      this.$nextTick(() => {
        const ta = document.querySelector('.composer-input');
        if (ta) {
          ta.style.height = 'auto';
          ta.style.height = Math.min(ta.scrollHeight, 120) + 'px';
        }
      });
    },

    async changePassword() {
      this.pwdError = '';
      if (this.pwdForm.new1 !== this.pwdForm.new2) {
        this.pwdError = 'رمزهای جدید یکسان نیستند';
        return;
      }
      try {
        await this.api('/api/auth/change-password', {
          method: 'POST',
          body: JSON.stringify({
            current_password: this.pwdForm.current,
            new_password: this.pwdForm.new1,
          }),
        });
        this.showChangePwd = false;
        this.pwdForm = { current: '', new1: '', new2: '' };
        if (this.session.user) this.session.user.must_change_pwd = false;
        alert('✅ رمز عبور تغییر کرد');
      } catch (e) {
        this.pwdError = e.message;
      }
    },

    /* ─────────── Admin: Users ─────────── */
    async createUser() {
      this.userFormError = '';
      try {
        const u = await this.api('/api/users', {
          method: 'POST',
          body: JSON.stringify(this.newUserForm),
        });
        this.users.push(u);
        this.showCreateUser = false;
        this.newUserForm = { username: '', display_name: '', password: '', role: 'user' };
        alert(`✅ کاربر ${u.display_name || u.username} ساخته شد`);
      } catch (e) {
        this.userFormError = e.message;
      }
    },

    async toggleRole(u) {
      const newRole = u.role === 'sup' ? 'user' : 'sup';
      try {
        await this.api(`/api/users/${u.id}/role`, {
          method: 'POST',
          body: JSON.stringify({ role: newRole }),
        });
        u.role = newRole;
      } catch (e) {
        alert(e.message);
      }
    },

    resetUserPwd(u) {
      this.resetPwdTarget = u;
      this.resetPwdForm = '';
      this.showResetPwd = true;
    },

    async confirmResetPwd() {
      if (!this.resetPwdForm) { alert('رمز جدید را وارد کنید'); return; }
      try {
        await this.api(`/api/users/${this.resetPwdTarget.id}/reset-password`, {
          method: 'POST',
          body: JSON.stringify({ new_password: this.resetPwdForm }),
        });
        this.showResetPwd = false;
        alert(`✅ رمز ${this.resetPwdTarget.display_name || this.resetPwdTarget.username} ریست شد`);
      } catch (e) {
        alert(e.message);
      }
    },

    async toggleUserActive(u) {
      try {
        await this.api(`/api/users/${u.id}/toggle-active`, { method: 'POST' });
        u.active = !u.active;
      } catch (e) {
        alert(e.message);
      }
    },

    async deleteUser(u) {
      if (!confirm(`آیا ${u.display_name || u.username} حذف شود؟ این عمل قابل بازگشت نیست.`)) return;
      try {
        await this.api(`/api/users/${u.id}`, { method: 'DELETE' });
        this.users = this.users.filter(x => x.id !== u.id);
        alert('🗑 کاربر حذف شد');
      } catch (e) {
        alert(e.message);
      }
    },

    /* ─────────── Admin: License ─────────── */
    async loadLicense() {
      try {
        const lic = await this.api('/api/settings/license');
        this.license = lic;
      } catch (e) {}
    },

    async applyLicense(event) {
      const file = event.target.files[0];
      if (!file) return;
      const form = new FormData();
      form.append('license', file);
      try {
        const res = await fetch('/api/settings/license', {
          method: 'POST',
          headers: { 'Authorization': 'Bearer ' + this.session.token },
          body: form,
        });
        const data = await res.json();
        this.license = data;
        if (data.applied) {
          alert(`✅ لایسنس معتبر! حداکثر کاربران: ${data.max_users}`);
        } else {
          alert(`⚠️ ${data.error || 'لایسنس نامعتبر'}\nمحدودیت به ${data.max_users} کاربر کاهش یافت.`);
        }
      } catch (e) {
        alert('خطا در اعمال لایسنس');
      }
    },

    async removeLicense() {
      if (!confirm('لایسنس حذف شود و به ۲۰ کاربر رایگان برگردیم؟')) return;
      try {
        const res = await this.api('/api/settings/license', { method: 'DELETE' });
        this.license = { state: 'free', max_users: 20, licensed_to: '', expiry: '', error: '' };
        alert('✅ لایسنس حذف شد — ۲۰ کاربر رایگان');
      } catch (e) {
        alert(e.message);
      }
    },

    /* ─────────── Admin: Network ─────────── */
    async loadNetwork() {
      try {
        const n = await this.api('/api/settings/network');
        this.networkForm = { address_type: n.address_type, ip: n.ip, dns: n.dns, port: n.port };
      } catch (e) {}
    },

    async saveNetwork() {
      try {
        await this.api('/api/settings/network', {
          method: 'POST',
          body: JSON.stringify(this.networkForm),
        });
        alert('✅ تنظیمات شبکه ذخیره شد\n(برای اعمال، سرور را ری‌استارت کنید)');
      } catch (e) {
        alert(e.message);
      }
    },
  },

  async mounted() {
    // Apply saved theme
    document.documentElement.setAttribute('data-theme', this.uiTheme);
    // Ask for notification permission (only when tab hidden or user interacts)
    this.requestNotifPermission();
    // Fetch server info
    try {
      this.serverInfo = await this.api('/api/settings/info');
    } catch (e) {}
    await this.checkSetup();
  },

  watch: {
    showAdminPanel(v) {
      if (v) {
        this.adminTab = 'users'; // ensure users tab is active
        this.loadLicense();
        this.loadNetwork();
      }
    },
  },
});

app.mount('#app');
