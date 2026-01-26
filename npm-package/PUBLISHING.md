# Vibepup npm Publishing (from any machine)

## 0) Test from build (before pushing to dev npm)

Run and test the package **from the built artifacts** before publishing.

**Option A – Run from repo (fast iteration)**  
Requires Go 1.22+ for TUI.

```bash
cd npm-package
npm run build          # build TUI binary (tui/vibepup-tui)
npm run run:local -- --tui     # test TUI
npm run run:local -- --watch   # test CLI engine
npm run run:local -- 3         # test N-iteration run
```

**Option B – Install from local tarball (simulate npm install)**  
Same as users installing from registry.

```bash
cd npm-package
npm run pack:local     # build TUI + npm pack → vibepup-1.0.2.tgz
npm install -g ./vibepup-1.0.2.tgz
vibepup --tui          # or vibepup --watch, vibepup 3, etc.
```

**Option C – `npm link` (global symlink for development)**  
Changes in the repo are reflected immediately when you run `vibepup`.

```bash
cd npm-package
npm run build
npm link
cd /any/project
vibepup --tui
```

---

## 1) Login
```bash
npm adduser
```

## 2) Verify package metadata
```bash
node -p "require('./package.json').name"
node -p "require('./package.json').version"
```

## 3) Publish
```bash
npm publish --access public
```

## Notes
- Publish must be run inside `npm-package/`.
- If you see a `bin` warning, ensure `package.json` has: `"bin": "bin/ralph.js"`.
- Windows users: WSL2 mode is recommended for publishing to ensure consistent environment behavior. Windows-native mode is supported but may lack some Linux-parity features.
