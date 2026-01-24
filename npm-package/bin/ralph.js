#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

const scriptPath = path.join(__dirname, '../lib/ralph.sh');
const allArgs = process.argv.slice(2);
const useTui = allArgs.includes('--tui');
const args = allArgs.filter((arg) => arg !== '--tui');

try {
  fs.chmodSync(scriptPath, '755');
} catch (err) {}

console.log('🐾 Vibepup is waking up...');

const shellOptions = { stdio: 'inherit', shell: false };

if (useTui) {
  const tuiDir = path.join(__dirname, '../tui');
  const binName = os.platform() === 'win32' ? 'vibepup-tui.exe' : 'vibepup-tui';
  const binPath = path.join(tuiDir, binName);

  if (fs.existsSync(binPath)) {
    const tui = spawn(binPath, args, shellOptions);
    tui.on('error', (err) => {
      console.error('❌ Failed to start Vibepup TUI.');
      console.error(String(err));
      process.exit(1);
    });
    tui.on('close', (code) => process.exit(code));
    return;
  }

  if (os.platform() !== 'win32' && fs.existsSync(tuiDir)) {
    const goCmd = spawn('go', ['run', '.', ...args], { ...shellOptions, cwd: tuiDir });
    goCmd.on('error', (err) => {
      console.error('❌ Failed to start Vibepup TUI.');
      console.error(String(err));
      process.exit(1);
    });
    goCmd.on('close', (code) => process.exit(code));
    return;
  }

  console.error('❌ Vibepup TUI not available.');
  console.error('   Build it with:');
  console.error('   cd ' + tuiDir + ' && go build -o ' + binName);
  process.exit(1);
}

let command = 'bash';
let cmdArgs = [scriptPath, ...args];
if (os.platform() === 'win32') {
  shellOptions.shell = true;
}

const vibepup = spawn(command, cmdArgs, shellOptions);

vibepup.on('error', (err) => {
    console.error('❌ Failed to start Vibepup.');
    if (err.code === 'ENOENT') {
        console.error('   Error: Bash not found. Please install Git Bash or WSL.');
    } else {
        console.error(String(err));
    }
    process.exit(1);
});

vibepup.on('close', (code) => {
    process.exit(code);
});
