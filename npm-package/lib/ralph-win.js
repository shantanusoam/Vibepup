#!/usr/bin/env node

const path = require('path');
const { spawnSync } = require('child_process');

const runnerPath = path.join(__dirname, 'runner', 'index.js');
const result = spawnSync(process.execPath, [runnerPath, ...process.argv.slice(2)], { stdio: 'inherit' });
process.exit(result.status || 0);
