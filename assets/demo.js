/* Design System Showcase Runtime — Sentinel DESIGN.md */

const appRoot = document.getElementById('app');
let currentPage = 'Main';
let currentParams = {};
const pageStates = {};
let appState = null;
const historyStack = [];

function createState(initial, target) {
  const listeners = {};
  const proxy = new Proxy(target, {
    get(obj, key) {
      if (key === '$on') return (key, fn) => { listeners[key] = listeners[key] || []; listeners[key].push(fn); };
      if (key === '$defaults') return (defaults) => {
        Object.entries(defaults).forEach(([k, v]) => { if (obj[k] === undefined) obj[k] = v; });
      };
      return obj[key];
    },
    set(obj, key, value) {
      const old = obj[key];
      obj[key] = value;
      if (old !== value && listeners[key]) listeners[key].forEach(fn => fn(value));
      return true;
    }
  });
  Object.assign(target, initial);
  return proxy;
}

function el(tag, className = '', attrs = {}) {
  const e = document.createElement(tag);
  if (className) e.className = className;
  Object.entries(attrs).forEach(([k, v]) => {
    if (k === 'children') {
      if (Array.isArray(v)) v.forEach(c => c && e.appendChild(c));
    } else if (k === 'text') {
      e.textContent = v;
    } else if (k.startsWith('on') && typeof v === 'function') {
      const ev = k.slice(2).toLowerCase();
      e.addEventListener(ev, v);
    } else if (k === 'disabled' || k === 'checked' || k === 'selected') {
      if (v) e.setAttribute(k, '');
    } else if (v !== undefined && v !== null) {
      e.setAttribute(k, v);
    }
  });
  return e;
}

function icon(name, size = 20) {
  const i = document.createElement('i');
  i.setAttribute('data-lucide', name);
  i.style.width = `${size}px`;
  i.style.height = `${size}px`;
  i.style.display = 'inline-block';
  return i;
}

function navigate(page, params = {}) {
  historyStack.push({ page: currentPage, params: currentParams });
  currentPage = page;
  currentParams = params;
  render();
}
function back() {
  const prev = historyStack.pop();
  if (prev) {
    currentPage = prev.page;
    currentParams = prev.params;
  }
  render();
}
window.demoNavigate = (page, params) => navigate(page, params);
window.demoCurrentPage = () => currentPage;

function render() {
  appRoot.innerHTML = '';
  const renderer = pages[currentPage];
  if (!renderer) { appRoot.appendChild(el('div', '', { text: `Unknown page: ${currentPage}` })); return; }
  if (!pageStates[currentPage]) pageStates[currentPage] = {};
  appState = createState(pageStates[currentPage], pageStates[currentPage]);
  appState.$defaults = (defs) => { Object.entries(defs).forEach(([k, v]) => { if (appState[k] === undefined) appState[k] = v; }); };
  appRoot.appendChild(renderer(currentParams));
  if (window.lucide) lucide.createIcons();
}

let toastContainer = null;
function ensureToastContainer() {
  if (!toastContainer) {
    toastContainer = el('div', 'toast-container');
    document.body.appendChild(toastContainer);
  }
  return toastContainer;
}
function showToast(message, variant = 'default') {
  const map = { default: 'info', success: 'check-circle', error: 'alert-circle', warning: 'alert-triangle', info: 'info' };
  const t = el('div', `toast ${variant}`, {
    children: [icon(map[variant] || 'info', 18), el('span', '', { text: message })]
  });
  ensureToastContainer().appendChild(t);
  requestAnimationFrame(() => t.classList.add('show'));
  setTimeout(() => { t.classList.remove('show'); setTimeout(() => t.remove(), 250); }, 3000);
}

function showDialog(message) {
  const scrim = el('div', 'overlay');
  const card = el('div', 'modal-card', {
    children: [
      el('div', 'modal-header', { text: 'Confirm' }),
      el('div', 'modal-body column gap-4', { children: [el('p', '', { text: message })] }),
      el('div', 'modal-footer', {
        children: [
          Button('OK', { variant: 'primary', onClick: () => { scrim.remove(); showToast('Confirmed'); } }),
          Button('Cancel', { variant: 'secondary', onClick: () => scrim.remove() })
        ]
      })
    ]
  });
  scrim.appendChild(card);
  document.body.appendChild(scrim);
  const closeOnEsc = (e) => { if (e.key === 'Escape') { scrim.remove(); document.removeEventListener('keydown', closeOnEsc); } };
  document.addEventListener('keydown', closeOnEsc);
  scrim.addEventListener('click', (e) => { if (e.target === scrim) scrim.remove(); });
}

// ── Component factories ─────────────────────────────────────────────────────

function Text(text, { weight } = {}) {
  const className = weight === 'horizontal' ? 'text-headline-md' : 'text-body-md';
  return el('span', className, { text });
}

function Button(label, { variant = 'primary', size, disabled = false, loading = false, className = '', onClick } = {}) {
  const base = `btn btn-${variant}${size ? ` btn-${size}` : ''} ${className}`;
  const attrs = { disabled };
  if (onClick && !disabled) attrs.onClick = onClick;
  if (loading) {
    attrs.children = [Spinner(16), el('span', '', { text: label })];
  } else {
    attrs.text = label;
  }
  return el('button', base, attrs);
}

function Badge(label, { variant = 'default' } = {}) {
  return el('span', `badge badge-${variant}`, { text: label });
}

function Card(title, children = []) {
  const c = el('div', 'card');
  if (title) c.appendChild(el('div', 'text-title-md text-primary mb-2', { text: title }));
  children.forEach(ch => c.appendChild(ch));
  return c;
}

function Input(placeholder, { value = '', onInput, type = 'text' } = {}) {
  const input = el('input', 'input', { type, placeholder, value });
  if (onInput) input.addEventListener('input', (e) => onInput(e.target.value));
  return el('div', 'input-wrapper', { children: [input] });
}

function Select(options, { onInput } = {}) {
  const select = el('select', 'select');
  options.forEach(opt => select.appendChild(el('option', '', { value: opt, text: opt })));
  if (onInput) select.addEventListener('change', (e) => onInput(e.target.value));
  return el('div', 'input-wrapper', { children: [select] });
}

function Textarea(placeholder) {
  return el('textarea', 'input', { placeholder, rows: 3 });
}

function Checkbox(label, { onClick } = {}) {
  const box = el('span', 'checkbox');
  const check = icon('check', 14);
  box.appendChild(check);
  let checked = false;
  const wrapper = el('label', 'checkbox-wrapper');
  const input = el('input', '', { type: 'checkbox', style: 'display:none' });
  wrapper.appendChild(input);
  wrapper.appendChild(box);
  wrapper.appendChild(el('span', 'text-body-md text-primary', { text: label }));
  wrapper.addEventListener('click', () => { checked = !checked; box.classList.toggle('checked', checked); if (onClick) onClick(); });
  return wrapper;
}

function Toggle(label, { onClick } = {}) {
  const track = el('span', 'switch');
  let checked = false;
  const wrapper = el('label', 'toggle-wrapper');
  const input = el('input', '', { type: 'checkbox', style: 'display:none' });
  wrapper.appendChild(input);
  wrapper.appendChild(track);
  wrapper.appendChild(el('span', 'text-body-md text-primary', { text: label }));
  wrapper.addEventListener('click', () => { checked = !checked; track.classList.toggle('checked', checked); if (onClick) onClick(); });
  return wrapper;
}

function Radio(label, { onClick } = {}) {
  const dot = el('span', 'radio-dot');
  const circle = el('span', 'radio');
  circle.appendChild(dot);
  const wrapper = el('label', 'checkbox-wrapper');
  const input = el('input', '', { type: 'radio', name: 'radio-group', style: 'display:none' });
  wrapper.appendChild(input);
  wrapper.appendChild(circle);
  wrapper.appendChild(el('span', 'text-body-md text-primary', { text: label }));
  wrapper.addEventListener('click', () => { document.querySelectorAll('.radio').forEach(r => r.classList.remove('checked')); circle.classList.add('checked'); if (onClick) onClick(); });
  return wrapper;
}

function Chip(label, { onClick } = {}) {
  return el('button', 'chip', { text: label, onClick });
}

function Avatar(initial) {
  return el('div', 'avatar', { text: initial });
}

function Table(headers, rows) {
  const table = el('table', 'table');
  const thead = el('thead');
  const trHead = el('tr');
  headers.forEach(h => trHead.appendChild(el('th', '', { text: h })));
  thead.appendChild(trHead);
  const tbody = el('tbody');
  rows.forEach(row => { const tr = el('tr'); row.forEach(cell => tr.appendChild(el('td', '', { text: cell }))); tbody.appendChild(tr); });
  table.appendChild(thead);
  table.appendChild(tbody);
  return el('div', 'table-container', { children: [table] });
}

function DataListItem(label, value) {
  return el('div', 'data-list-item', {
    children: [
      el('div', 'data-list-label', { text: label }),
      el('div', 'data-list-value', { text: value })
    ]
  });
}

function EmptyState(title, children = []) {
  const es = el('div', 'empty-state column align-center gap-3');
  es.appendChild(icon('inbox', 48));
  es.appendChild(el('div', 'empty-state-title', { text: title }));
  children.forEach(ch => es.appendChild(ch));
  return es;
}

function Alert(message, { variant = 'default' } = {}) {
  const map = { default: 'info', info: 'info', success: 'check-circle', warning: 'alert-triangle', error: 'alert-circle' };
  return el('div', `alert alert-${variant}`, {
    children: [
      icon(map[variant] || 'info', 20),
      el('div', 'column gap-1', { children: [el('div', 'alert-title', { text: variant.charAt(0).toUpperCase() + variant.slice(1) }), el('div', 'alert-description', { text: message })] })
    ]
  });
}

function Progress(value, { max = 100 } = {}) {
  const fill = el('div', 'progress-fill');
  fill.style.width = `${Math.min(100, Math.max(0, (value / max) * 100))}%`;
  return el('div', 'progress', { children: [fill] });
}

function Skeleton() {
  return el('div', 'skeleton', { style: 'height: 16px; margin-bottom: 8px;' });
}

function Spinner(size = 20) {
  return icon('loader-2', size);
}

// ── Layout helpers ──────────────────────────────────────────────────────────

function row(children, { wrap = false, className = '' } = {}) {
  return el('div', `row${wrap ? ' wrap' : ''} ${className}`.trim(), { children });
}
function column(children, { className = '' } = {}) {
  return el('div', `column ${className}`.trim(), { children });
}
function SectionHeader(title, description) {
  return column([
    el('h1', 'text-display-sm text-primary', { text: title }),
    el('p', 'text-body-md text-secondary', { text: description }),
    el('div', 'divider my-3')
  ]);
}
function DemoCard(title, children) {
  return Card(title, children);
}
function NavButton(label, target) {
  return Button(label, { variant: 'ghost', onClick: () => navigate(typeof target === 'string' ? target : 'Main') });
}

function ColorSwatch(name) {
  return column([
    el('div', 'color-swatch'),
    el('div', 'text-caption-md text-secondary mt-1', { text: name })
  ]);
}
function TypeSample(label, size) {
  return row([
    el('span', 'text-body-md text-primary', { text: label }),
    el('span', 'text-caption-md text-secondary', { text: size })
  ], { className: 'align-center justify-between' });
}
function SpacingBlock(label) {
  return column([
    el('div', 'spacing-block'),
    el('div', 'text-caption-md text-secondary mt-1', { text: label })
  ]);
}

// ── Pages ───────────────────────────────────────────────────────────────────

const pages = {};

pages.Main = () => {
  const s = appState;
  s.$defaults({ activeSection: 'Colors' });
  return column([
    column([
      Badge('DESIGN.md Showcase', { variant: 'info' }),
      el('h1', 'hero-title', { text: 'Sentinel Design System' }),
      el('p', 'hero-subtitle', { text: 'Explore tokens, components, and patterns defined in the design system file.' }),
      Button('Browse Tokens', { onClick: () => navigate('Colors') })
    ], { className: 'hero column align-center' }),
    row([
      Card('Colors', [el('div', 'text-headline-md text-primary', { text: '50+' }), el('div', 'text-body-sm text-secondary', { text: 'Semantic roles' })]),
      Card('Components', [el('div', 'text-headline-md text-primary', { text: '75+' }), el('div', 'text-body-sm text-secondary', { text: 'Named variants' })]),
      Card('Type Styles', [el('div', 'text-headline-md text-primary', { text: '24' }), el('div', 'text-body-sm text-secondary', { text: 'Scale levels' })])
    ], { wrap: true, className: 'stats-ribbon' }),
    row([
      Button('Colors', { variant: 'secondary', onClick: () => navigate('Colors') }),
      Button('Typography', { variant: 'secondary', onClick: () => navigate('Typography') }),
      Button('Spacing', { variant: 'secondary', onClick: () => navigate('Spacing') }),
      Button('Components', { variant: 'secondary', onClick: () => navigate('Components') }),
      Button('Layouts', { variant: 'secondary', onClick: () => navigate('Layouts') }),
      Button('Feedback', { variant: 'secondary', onClick: () => navigate('Feedback') }),
      Button('Accessibility', { variant: 'secondary', onClick: () => navigate('Accessibility') })
    ], { wrap: true, className: 'toolbar' }),
    DemoCard('System Principles', [
      column([
        el('div', 'text-body-md text-primary', { text: 'Information density without clutter' }),
        el('div', 'text-body-md text-primary', { text: 'Accessibility as a feature' }),
        el('div', 'text-body-md text-primary', { text: 'Tokens are the single source of truth' }),
        el('div', 'text-body-md text-primary', { text: 'Dark mode and reduced motion ready' })
      ])
    ])
  ], { className: 'column gap-6 w-full' });
};

pages.Colors = () => {
  return column([
    SectionHeader('Colors', 'Semantic color roles for surfaces, text, outlines, and status'),
    NavButton('← Back', 'Main'),
    DemoCard('Brand', [row([ColorSwatch('Primary'), ColorSwatch('Secondary'), ColorSwatch('Tertiary')], { wrap: true })]),
    DemoCard('Text & Icons', [row([ColorSwatch('Text Primary'), ColorSwatch('Text Secondary'), ColorSwatch('Text Tertiary'), ColorSwatch('Text Inverse'), ColorSwatch('Text Link')], { wrap: true })]),
    DemoCard('Surfaces', [row([ColorSwatch('Surface'), ColorSwatch('Surface Elevated'), ColorSwatch('Neutral'), ColorSwatch('Surface Inverse')], { wrap: true })]),
    DemoCard('Status', [row([ColorSwatch('Success'), ColorSwatch('Warning'), ColorSwatch('Error'), ColorSwatch('Info')], { wrap: true })]),
    DemoCard('Outlines', [row([ColorSwatch('Outline'), ColorSwatch('Outline Hover'), ColorSwatch('Outline Focus'), ColorSwatch('Outline Error')], { wrap: true })])
  ], { className: 'column gap-4 w-full' });
};

pages.Typography = () => {
  return column([
    SectionHeader('Typography', 'Display, headline, title, body, label, caption, code, and overline styles'),
    NavButton('← Back', 'Main'),
    DemoCard('Display & Headlines', [
      column([
        TypeSample('Display 2XL', '72px / 800'),
        TypeSample('Display XL', '60px / 800'),
        TypeSample('Headline LG', '24px / 600'),
        TypeSample('Headline MD', '20px / 600')
      ])
    ]),
    DemoCard('Body & Captions', [
      column([
        TypeSample('Body LG', '18px / 400'),
        TypeSample('Body MD', '16px / 400'),
        TypeSample('Body SM', '14px / 400'),
        TypeSample('Caption MD', '12px / 400')
      ])
    ]),
    DemoCard('Labels & Code', [
      column([
        TypeSample('Label LG', '14px / 500 Mono'),
        TypeSample('Label MD', '12px / 500 Mono'),
        TypeSample('Code MD', '13px / 400 Mono'),
        TypeSample('Overline', '10px / 700')
      ])
    ])
  ], { className: 'column gap-4 w-full' });
};

pages.Spacing = () => {
  return column([
    SectionHeader('Spacing', '4px baseline scale, layout breakpoints, icon sizes, and hit targets'),
    NavButton('← Back', 'Main'),
    DemoCard('Base Scale', [row([SpacingBlock('0'), SpacingBlock('1 (4px)'), SpacingBlock('2 (8px)'), SpacingBlock('3 (12px)'), SpacingBlock('4 (16px)'), SpacingBlock('6 (24px)'), SpacingBlock('8 (32px)')], { wrap: true })]),
    DemoCard('Layout Breakpoints', [row([SpacingBlock('XS 480px'), SpacingBlock('SM 640px'), SpacingBlock('MD 768px'), SpacingBlock('LG 1024px'), SpacingBlock('XL 1280px')], { wrap: true })]),
    DemoCard('Icon Scale', [row([SpacingBlock('XS 12px'), SpacingBlock('SM 16px'), SpacingBlock('MD 20px'), SpacingBlock('LG 24px'), SpacingBlock('XL 32px')], { wrap: true })])
  ], { className: 'column gap-4 w-full' });
};

pages.Components = () => {
  const s = appState;
  s.$defaults({ pressedCount: 0 });
  return column([
    SectionHeader('Components', 'Buttons, inputs, selects, toggles, badges, chips, and avatars'),
    NavButton('← Back', 'Main'),
    DemoCard('Buttons', [
      row([
        Button('Primary', { onClick: () => s.pressedCount++ }),
        Button('Secondary', { variant: 'secondary', onClick: () => s.pressedCount++ }),
        Button('Tertiary', { variant: 'tertiary', onClick: () => s.pressedCount++ }),
        Button('Ghost', { variant: 'ghost', onClick: () => s.pressedCount++ }),
        Button('Danger', { variant: 'danger', onClick: () => showDialog('Confirm destructive action?') }),
        Button('Disabled', { disabled: true })
      ], { wrap: true })
    ]),
    DemoCard('Inputs & Selects', [
      column([
        Input('Email address', { placeholder: 'you@example.com' }),
        Select(['United States', 'Canada', 'United Kingdom', 'Germany', 'Japan']),
        Textarea('Notes')
      ])
    ]),
    DemoCard('Toggles', [
      row([
        Checkbox('I agree to terms', { onClick: () => showToast('Agreement toggled') }),
        Toggle('Enable notifications', { onClick: () => showToast('Notifications toggled') }),
        Radio('Standard plan', { onClick: () => showToast('Plan selected') })
      ], { wrap: true })
    ]),
    DemoCard('Badges, Chips & Tags', [
      row([
        Badge('Active', { variant: 'success' }),
        Badge('Warning', { variant: 'warning' }),
        Badge('Error', { variant: 'error' }),
        Badge('Info', { variant: 'info' }),
        Chip('Filter A', { onClick: () => showToast('Filter toggled') }),
        Chip('Filter B', { onClick: () => showToast('Filter toggled') })
      ], { wrap: true })
    ]),
    DemoCard('Avatars', [
      row([Avatar('A'), Avatar('B'), Avatar('C'), Avatar('D')], { wrap: true })
    ])
  ], { className: 'column gap-4 w-full' });
};

pages.Layouts = () => {
  return column([
    SectionHeader('Layouts', 'Cards, tables, data lists, empty states, and responsive grids'),
    NavButton('← Back', 'Main'),
    DemoCard('Cards', [
      row([
        Card('Card', [el('div', 'text-body-sm text-secondary', { text: 'Default surface card' })]),
        Card('Elevated', [el('div', 'text-body-sm text-secondary', { text: 'Elevated surface variant' })]),
        Card('Bordered', [el('div', 'text-body-sm text-secondary', { text: 'Outlined card variant' })])
      ], { wrap: true })
    ]),
    DemoCard('Table', [
      Table(['Name', 'Role', 'Status'], [
        ['Alice Johnson', 'Admin', 'Active'],
        ['Bob Smith', 'Developer', 'Active'],
        ['Carol Williams', 'Designer', 'Away']
      ])
    ]),
    DemoCard('Data List', [
      column([
        DataListItem('Plan', 'Professional'),
        DataListItem('Seats', '12 of 20 used'),
        DataListItem('Renewal', '2024-12-15')
      ])
    ]),
    DemoCard('Empty State', [
      EmptyState('No results found', [
        el('div', 'text-body-md text-secondary', { text: 'Try adjusting your filters or search query.' }),
        Button('Clear Filters', { variant: 'secondary', onClick: () => showToast('Filters cleared') })
      ])
    ])
  ], { className: 'column gap-4 w-full' });
};

pages.Feedback = () => {
  const s = appState;
  s.$defaults({ progress: 35, showSkeleton: false });
  return column([
    SectionHeader('Feedback', 'Alerts, toasts, spinners, progress bars, and skeleton states'),
    NavButton('← Back', 'Main'),
    DemoCard('Alerts', [
      column([
        Alert('Everything is running smoothly.', { variant: 'info' }),
        Alert('Your changes were saved successfully.', { variant: 'success' }),
        Alert('Subscription expires in 3 days.', { variant: 'warning' }),
        Alert('Unable to connect to the server.', { variant: 'error' })
      ])
    ]),
    DemoCard('Toasts', [
      row([
        Button('Show Success', { variant: 'secondary', onClick: () => showToast('Operation completed', 'success') }),
        Button('Show Error', { variant: 'secondary', onClick: () => showToast('Something went wrong', 'error') }),
        Button('Show Warning', { variant: 'secondary', onClick: () => showToast('Proceed with caution', 'warning') })
      ], { wrap: true })
    ]),
    DemoCard('Progress', [
      column([
        el('div', 'text-body-md text-primary', { text: `Task completion: ${s.progress}%` }),
        Progress(s.progress),
        row([
          Button('-10', { variant: 'secondary', onClick: () => { if (s.progress >= 10) s.progress -= 10; } }),
          Button('+10', { variant: 'secondary', onClick: () => { if (s.progress < 100) s.progress += 10; } }),
          Button('Reset', { variant: 'ghost', onClick: () => s.progress = 0 })
        ], { wrap: true })
      ])
    ]),
    DemoCard('Loading States', [
      column([
        Button(s.showSkeleton ? 'Hide Skeleton' : 'Show Skeleton', { variant: 'secondary', onClick: () => s.showSkeleton = !s.showSkeleton }),
        s.showSkeleton ? column([Skeleton(), Skeleton(), Skeleton()]) : el('span')
      ])
    ])
  ], { className: 'column gap-4 w-full' });
};

pages.Accessibility = () => {
  const s = appState;
  s.$defaults({ reducedMotion: false, highContrast: false });
  return column([
    SectionHeader('Accessibility', 'Focus rings, keyboard navigation, and preference toggles'),
    NavButton('← Back', 'Main'),
    DemoCard('Focus Order', [
      row([
        Input('First name'),
        Input('Last name'),
        Button('Submit', { onClick: () => showToast('Submitted', 'success') })
      ], { wrap: true })
    ]),
    DemoCard('Preference Toggles', [
      column([
        Toggle('Reduced Motion', { onClick: () => { s.reducedMotion = !s.reducedMotion; showToast(s.reducedMotion ? 'Reduced motion on' : 'Reduced motion off'); } }),
        Toggle('High Contrast', { onClick: () => { s.highContrast = !s.highContrast; showToast(s.highContrast ? 'High contrast on' : 'High contrast off'); } })
      ])
    ]),
    DemoCard('Keyboard Shortcuts', [
      column([
        el('div', 'text-body-md text-primary', { text: 'Ctrl+K — Open command palette' }),
        el('div', 'text-body-md text-primary', { text: 'Ctrl+S — Save' }),
        el('div', 'text-body-md text-primary', { text: 'Esc — Close modal or dialog' }),
        el('div', 'text-body-md text-primary', { text: 'Tab — Move to next focusable element' })
      ])
    ])
  ], { className: 'column gap-4 w-full' });
};

// ── Boot ────────────────────────────────────────────────────────────────────

render();
