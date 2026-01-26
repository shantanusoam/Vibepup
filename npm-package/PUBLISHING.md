# Vibepup npm Publishing (from any machine)

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
