(() => {
  const preference = matchMedia('(prefers-color-scheme: dark)');
  const applyTheme = () => {
    let saved;
    try { saved = localStorage.getItem('ticket-theme'); } catch {}
    document.documentElement.dataset.theme = ['light', 'dark'].includes(saved) ? saved : (preference.matches ? 'dark' : 'light');
    syncThemeControls();
  };
  function syncThemeControls() {
    const dark = document.documentElement.dataset.theme === 'dark';
    const label = dark ? 'Switch to light mode' : 'Switch to dark mode';
    document.querySelectorAll('[data-action=theme]').forEach(control => {
      control.setAttribute('aria-label', label);
      control.setAttribute('title', label);
      const text = control.querySelector('[data-theme-label]');
      if (text) text.textContent = label;
    });
  }
  applyTheme();
  preference.addEventListener('change', applyTheme);
  const mobile = matchMedia('(max-width: 767px)');
  let navigationOpener;
  let accountOpener;
  try {
    document.documentElement.classList.toggle('sidebar-collapsed', localStorage.getItem('support-sidebar') === 'collapsed');
  } catch {}

  function syncNavigation() {
    const sidebar = document.querySelector('#sidebar');
    const open = mobile.matches && document.documentElement.classList.contains('sidebar-open');
    if (sidebar) {
      sidebar.inert = mobile.matches && !open;
      if (open) {
        sidebar.setAttribute('role', 'dialog');
        sidebar.setAttribute('aria-modal', 'true');
      } else {
        sidebar.removeAttribute('role');
        sidebar.removeAttribute('aria-modal');
      }
    }
    const body = document.querySelector('.workspace-body');
    if (body) body.inert = open;
    const backdrop = document.querySelector('.sidebar-backdrop');
    if (backdrop) backdrop.hidden = !open;
    const toggle = document.querySelector('.sidebar-toggle');
    if (toggle) {
      const expanded = mobile.matches ? open : !document.documentElement.classList.contains('sidebar-collapsed');
      const label = mobile.matches ? (open ? 'Close navigation' : 'Open navigation') : (expanded ? 'Collapse sidebar' : 'Expand sidebar');
      toggle.setAttribute('aria-expanded', String(expanded));
      toggle.setAttribute('aria-label', label);
      toggle.setAttribute('title', label);
      toggle.querySelector('use')?.setAttribute('href', `/static/icons.svg#${mobile.matches ? 'menu' : expanded ? 'panel-close' : 'panel-open'}`);
    }
  }

  function closeNavigation(restoreFocus = true) {
    const wasOpen = document.documentElement.classList.contains('sidebar-open');
    document.documentElement.classList.remove('sidebar-open');
    syncNavigation();
    if (wasOpen && restoreFocus) navigationOpener?.focus();
  }

  function initializeWorkspace() {
    closeNavigation(false);
    syncThemeControls();
    const now = new Date();
    document.querySelectorAll('[data-current-year]').forEach(el => { el.textContent = now.getFullYear(); });
  }
  function bindAccountDialog() {
    const dialog = document.querySelector('#account-dialog');
    if (!dialog || dialog.dataset.bound) return;
    dialog.dataset.bound = 'true';
    dialog.addEventListener('close', () => accountOpener?.focus());
    dialog.addEventListener('click', event => {
      if (event.target !== dialog) return;
      const rect = dialog.getBoundingClientRect();
      if (event.clientX < rect.left || event.clientX > rect.right || event.clientY < rect.top || event.clientY > rect.bottom) dialog.close();
    });
  }
  bindAccountDialog();
  initializeWorkspace();
  mobile.addEventListener('change', () => closeNavigation());

  document.addEventListener('keydown', event => {
    if (event.key === 'Escape' && !document.querySelector('#account-dialog[open]')) {
      closeNavigation();
      const account = document.querySelector('.account-dropdown[open]');
      if (account) { account.open = false; account.querySelector('summary')?.focus(); }
    }
    if (event.key !== 'Tab' || !mobile.matches || !document.documentElement.classList.contains('sidebar-open')) return;
    if (document.querySelector('#account-dialog[open]')) return;
    const controls = [...document.querySelectorAll('#sidebar a, #sidebar button, #sidebar summary')].filter(el => el.getClientRects().length);
    const first = controls[0], last = controls.at(-1);
    if (event.shiftKey && document.activeElement === first) { event.preventDefault(); last?.focus(); }
    if (!event.shiftKey && document.activeElement === last) { event.preventDefault(); first?.focus(); }
  });

  document.addEventListener('click', event => {
    const account = document.querySelector('.account-dropdown[open]');
    if (account && !account.contains(event.target)) account.open = false;
    const anchor = event.target.closest('a[href^="#"]');
    if (anchor) {
      closeNavigation(false);
      if (account) account.open = false;
      const target = document.getElementById(anchor.getAttribute('href').slice(1));
      if (target) {
        target.setAttribute('tabindex', '-1');
        target.focus({ preventScroll: true });
      }
    }
    const control = event.target.closest('[data-action]');
    if (!control) return;
    if (control.dataset.action === 'account') {
      closeNavigation(false);
      if (account) account.open = false;
      accountOpener = mobile.matches ? navigationOpener : document.querySelector('.account-dropdown summary');
      document.querySelector('#account-dialog')?.showModal();
    }
    if (control.dataset.action === 'close-account') document.querySelector('#account-dialog')?.close();
    if (control.dataset.action === 'close-sidebar') closeNavigation();
    if (control.dataset.action === 'sidebar') {
      if (mobile.matches) {
        navigationOpener = control;
        document.documentElement.classList.toggle('sidebar-open');
        syncNavigation();
        document.querySelector('#sidebar .sidebar-close')?.focus();
      } else {
        const collapsed = document.documentElement.classList.toggle('sidebar-collapsed');
        try { localStorage.setItem('support-sidebar', collapsed ? 'collapsed' : 'expanded'); } catch {}
        syncNavigation();
      }
    }
    if (control.dataset.action === 'theme') {
      const theme = document.documentElement.dataset.theme === 'dark' ? 'light' : 'dark';
      document.documentElement.dataset.theme = theme;
      syncThemeControls();
      try { localStorage.setItem('ticket-theme', theme); } catch {}
    }
    if (control.dataset.action === 'toggle-password') {
      const input = document.getElementById(control.getAttribute('aria-controls'));
      if (!input) return;
      const visible = input.type === 'password';
      input.type = visible ? 'text' : 'password';
      control.textContent = visible ? 'Hide' : 'Show';
      control.setAttribute('aria-label', visible ? 'Hide password' : 'Show password');
      control.setAttribute('aria-pressed', String(visible));
    }
  });
  document.addEventListener('htmx:beforeSwap', event => {
    if ([400, 401, 409, 422, 429, 503].includes(event.detail.xhr.status) && event.detail.xhr.getResponseHeader('Content-Type')?.includes('text/html')) {
      event.detail.shouldSwap = true;
      event.detail.isError = false;
    }
  });
  document.addEventListener('htmx:afterSwap', () => {
    initializeWorkspace();
    bindAccountDialog();
    document.querySelector('#network-error')?.setAttribute('hidden', '');
    document.querySelector('[role=alert][tabindex], #main')?.focus({preventScroll: true});
    const heading = document.querySelector('#main h1, .auth-form h1');
    if (heading) document.title = `${heading.textContent.trim()} · IT Support`;
  });
  const showError = () => document.querySelector('#network-error')?.removeAttribute('hidden');
  document.addEventListener('htmx:sendError', showError);
  document.addEventListener('htmx:responseError', showError);
  window.addEventListener('pageshow', event => { if (event.persisted) location.reload(); });
})();
