#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

const scriptPath = path.join(__dirname, '../lib/ralph.sh');
const args = process.argv.slice(2);

try {
  fs.chmodSync(scriptPath, '755');
} catch (err) {}

console.log('🐾 Vibepup is waking up...');

let command = scriptPath;
let cmdArgs = args;
let shellOptions = { stdio: 'inherit', shell: false };

if (os.platform() === 'win32') {
    command = 'bash';
    cmdArgs = [scriptPath, ...args];
    shellOptions.shell = true;
} else {
    command = 'bash';
    cmdArgs = [scriptPath, ...args];
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
