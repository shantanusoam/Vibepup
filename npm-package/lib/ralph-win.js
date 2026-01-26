#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const crypto = require('crypto');
const { spawn, spawnSync } = require('child_process');

const ENGINE_DIR = path.resolve(__dirname);
const PROJECT_DIR = process.cwd();
const RUNS_DIR = path.join(PROJECT_DIR, '.ralph', 'runs');

const DEFAULT_ITERATIONS = 5;
const RALPH_MAX_TURN_SECONDS = Number.parseInt(process.env.RALPH_MAX_TURN_SECONDS || '900', 10);
const RALPH_NO_OUTPUT_SECONDS = Number.parseInt(process.env.RALPH_NO_OUTPUT_SECONDS || '180', 10);

const BUILD_MODELS_PREF = [
  'github-copilot/gpt-5.2-codex',
  'github-copilot/claude-sonnet-4.5',
  'github-copilot/gemini-3-pro-preview',
  'github-copilot-enterprise/gpt-5.2-codex',
  'github-copilot-enterprise/claude-sonnet-4.5',
  'github-copilot-enterprise/gemini-3-pro-preview',
  'openai/gpt-5.2-codex',
  'openai/gpt-5.1-codex-max',
  'google/gemini-3-pro-preview',
  'opencode/grok-code',
];

const PLAN_MODELS_PREF = [
  'github-copilot/claude-opus-4.5',
  'github-copilot/gemini-3-pro-preview',
  'github-copilot-enterprise/claude-opus-4.5',
  'github-copilot-enterprise/gemini-3-pro-preview',
  'openai/gpt-5.2',
  'google/antigravity-claude-opus-4-5-thinking',
  'google/gemini-3-pro-preview',
  'opencode/glm-4.7-free',
];

const SYSTEM_PROMPT = path.join(ENGINE_DIR, 'prompt.md');
const ARCHITECT_FILE = path.join(ENGINE_DIR, 'agents', 'architect.md');

const isWindows = process.platform === 'win32';

const parseArgs = () => {
  const args = process.argv.slice(2);
  let iterations = DEFAULT_ITERATIONS;
  let watchMode = false;
  let mode = 'default';
  let projectIdea = '';
  let freeMode = false;

  let index = 0;
  while (index < args.length) {
    const arg = args[index];
    if (arg === 'free') {
      freeMode = true;
      index += 1;
      continue;
    }
    if (arg === 'new') {
      mode = 'new';
      projectIdea = args[index + 1] || '';
      index += 2;
      continue;
    }
    if (arg === '--watch') {
      watchMode = true;
      index += 1;
      continue;
    }
    if (/^\d+$/.test(arg)) {
      iterations = Number.parseInt(arg, 10);
      index += 1;
      continue;
    }
    index += 1;
  }

  return {
    iterations,
    watchMode,
    mode,
    projectIdea,
    freeMode,
  };
};

const ensureDir = (dir) => fs.mkdirSync(dir, { recursive: true });

const md5File = (filePath) => {
  const content = fs.readFileSync(filePath, 'utf8');
  return crypto.createHash('md5').update(content).digest('hex');
};

const fileExists = (filePath) => fs.existsSync(filePath);

const readTail = (filePath, maxLines) => {
  if (!fileExists(filePath)) return '';
  const content = fs.readFileSync(filePath, 'utf8');
  const lines = content.split(/\r?\n/);
  return lines.slice(Math.max(0, lines.length - maxLines)).join('\n');
};

const ensureProjectFiles = () => {
  if (!fileExists(path.join(PROJECT_DIR, 'prd.md'))) {
    if (fileExists(path.join(PROJECT_DIR, 'prd.json'))) {
      console.log('🔄 Migrating legacy prd.json to prd.md...');
      const data = JSON.parse(fs.readFileSync(path.join(PROJECT_DIR, 'prd.json'), 'utf8'));
      const lines = data.map((item) => `- [ ] ${item.description}`);
      fs.writeFileSync(path.join(PROJECT_DIR, 'prd.md'), lines.join('\n') + '\n', 'utf8');
      fs.renameSync(path.join(PROJECT_DIR, 'prd.json'), path.join(PROJECT_DIR, 'prd.json.bak'));
    } else {
      console.log('⚠️  No prd.md found. Initializing...');
      const init = [
        '# Product Requirements Document (PRD)',
        '',
        '- [ ] Initialize repo-map.md with project architecture',
        '- [ ] Setup initial project structure',
        '',
      ].join('\n');
      fs.writeFileSync(path.join(PROJECT_DIR, 'prd.md'), init, 'utf8');
    }
  }

  if (!fileExists(path.join(PROJECT_DIR, 'repo-map.md'))) {
    fs.writeFileSync(path.join(PROJECT_DIR, 'repo-map.md'), '', 'utf8');
  }

  if (!fileExists(path.join(PROJECT_DIR, 'prd.state.json'))) {
    fs.writeFileSync(path.join(PROJECT_DIR, 'prd.state.json'), '{}', 'utf8');
  }

  if (!fileExists(path.join(PROJECT_DIR, 'progress.log'))) {
    fs.writeFileSync(path.join(PROJECT_DIR, 'progress.log'), '', 'utf8');
  }
};

const detectPhase = () => {
  const repoMapPath = path.join(PROJECT_DIR, 'repo-map.md');
  if (!fileExists(repoMapPath)) return 'PLAN';
  const content = fs.readFileSync(repoMapPath, 'utf8');
  return content.trim().length === 0 ? 'PLAN' : 'BUILD';
};

const resolveAvailableModels = (prefModels) => {
  console.error('🔍 Verifying available models...');
  const result = spawnSync('opencode', ['models', '--refresh'], { encoding: 'utf8' });
  const output = result.stdout || '';
  const lines = output.split(/\r?\n/).filter((line) => /^[a-z0-9-]+\/[a-z0-9.-]+$/.test(line));
  const available = [];
  for (const pref of prefModels) {
    if (lines.includes(pref)) {
      available.push(pref);
    }
  }
  if (available.length === 0) {
    console.error('⚠️  No preferred models found. Falling back to generic discovery.');
    const gptFallback = lines.find((line) => line.includes('gpt-4o'));
    const claudeFallback = lines.find((line) => line.includes('claude-sonnet'));
    if (gptFallback) available.push(gptFallback);
    if (claudeFallback) available.push(claudeFallback);
  }
  if (available.length === 0) {
    available.push('opencode/grok-code');
    console.error('⚠️  Using fallback model: opencode/grok-code');
  }
  return available;
};

const runWithWatchdog = (logPath, command, args) => new Promise((resolve) => {
  fs.writeFileSync(logPath, '', 'utf8');
  const logStream = fs.createWriteStream(logPath, { flags: 'a' });
  const child = spawn(command, args, { stdio: ['ignore', 'pipe', 'pipe'] });
  let lastOutput = Date.now();
  let killed = false;

  const handleData = (data) => {
    lastOutput = Date.now();
    logStream.write(data);
    process.stdout.write(data);
  };

  child.stdout.on('data', handleData);
  child.stderr.on('data', handleData);

  const interval = setInterval(() => {
    const now = Date.now();
    if (now - lastOutput > RALPH_NO_OUTPUT_SECONDS * 1000) {
      logStream.write('[RALPH] NO OUTPUT: likely waiting for input / hung tool\n');
      if (!killed) {
        killed = true;
        child.kill('SIGINT');
        setTimeout(() => child.kill('SIGTERM'), 3000);
        setTimeout(() => child.kill('SIGKILL'), 4000);
      }
    }
    if (now - startTime > RALPH_MAX_TURN_SECONDS * 1000) {
      logStream.write('[RALPH] TIMEOUT: killing opencode turn\n');
      if (!killed) {
        killed = true;
        child.kill('SIGINT');
        setTimeout(() => child.kill('SIGTERM'), 3000);
        setTimeout(() => child.kill('SIGKILL'), 4000);
      }
    }
  }, 5000);

  const startTime = Date.now();

  child.on('close', (code) => {
    clearInterval(interval);
    logStream.end();
    resolve(code || 0);
  });
});

const runAgent = async (model, phase, iterDir) => {
  const logPath = path.join(iterDir, 'agent_response.txt');
  const promptSuffix = phase === 'PLAN'
    ? 'MODE: PLAN. Focus on exploring and mapping. Do NOT write implementation code yet.'
    : 'MODE: BUILD. Focus on completing tasks in prd.md.';

  const args = [
    'run',
    `Proceed with task. ${promptSuffix}`,
    '--file', SYSTEM_PROMPT,
    '--file', path.join(PROJECT_DIR, 'prd.md'),
    '--file', path.join(PROJECT_DIR, 'prd.state.json'),
    '--file', path.join(PROJECT_DIR, 'repo-map.md'),
    '--file', path.join(iterDir, 'progress.tail.log'),
    '--model', model,
  ];

  return runWithWatchdog(logPath, 'opencode', args);
};

const runArchitect = () => {
  console.log('🏗️  Phase 0: The Architect');
  const args = [
    'run',
    `PROJECT IDEA: ${projectIdea}`,
    '--file', ARCHITECT_FILE,
    '--agent', 'general',
    '--model', planModels[0],
  ];
  const result = spawnSync('opencode', args, { stdio: 'inherit' });
  return result.status || 0;
};

const ensureOpencode = (freeMode) => {
  const exists = spawnSync('opencode', ['--version'], { stdio: 'ignore' }).status === 0;
  if (exists) return true;

  if (freeMode) {
    console.log('🔧 Free setup: installing opencode...');
    const npmAvailable = spawnSync('npm', ['--version'], { stdio: 'ignore' }).status === 0;
    if (!npmAvailable) {
      console.error('❌ npm not found. Install Node.js or use WSL2 for full setup.');
      return false;
    }
    spawnSync('npm', ['install', '-g', 'opencode-ai', 'opencode-antigravity-auth'], { stdio: 'inherit' });
  } else {
    console.error('❌ opencode not found. Vibepup requires opencode to run.');
    console.error('   Install with: npm install -g opencode-ai');
    console.error('   Free-tier option: vibepup free');
    return false;
  }
  return true;
};

const runFreeSetup = () => {
  console.log('✨ Vibepup Free Setup');
  console.log('   1) Installing auth plugin');
  spawnSync('npm', ['install', '-g', 'opencode-antigravity-auth'], { stdio: 'inherit' });
  console.log('   2) Starting Google auth');
  spawnSync('opencode', ['auth', 'login', 'antigravity'], { stdio: 'inherit' });
  console.log('   3) Refreshing models');
  spawnSync('opencode', ['models', '--refresh'], { stdio: 'inherit' });
  console.log("✅ Free setup complete. Run 'vibepup --watch' next.");
  process.exit(0);
};

const { iterations, watchMode, mode, projectIdea, freeMode } = parseArgs();

console.log('🐾 Vibepup v1.0 (Windows Native CLI Mode)');
console.log(`   Engine:  ${ENGINE_DIR}`);
console.log(`   Context: ${PROJECT_DIR}`);
console.log('   Tips:');
console.log("   - Run 'vibepup free' for free-tier setup");
console.log("   - Run 'vibepup new \"My idea\"' to bootstrap a project");
console.log("   - Run 'vibepup --tui' for a guided interface");

ensureDir(RUNS_DIR);
ensureProjectFiles();

if (!ensureOpencode(freeMode)) {
  process.exit(127);
}

if (freeMode) {
  runFreeSetup();
}

const buildModels = resolveAvailableModels(BUILD_MODELS_PREF);
const planModels = resolveAvailableModels(PLAN_MODELS_PREF);

if (mode === 'new') {
  const code = runArchitect();
  if (code !== 0) process.exit(code);
  console.log('✅ Architect initialization complete.');
}

let lastHash = md5File(path.join(PROJECT_DIR, 'prd.md'));
let i = 1;

const runLoop = async () => {
  while (true) {
    const currentHash = md5File(path.join(PROJECT_DIR, 'prd.md'));
    if (currentHash !== lastHash) {
      console.log('👀 PRD Changed! Restarting loop...');
      fs.appendFileSync(path.join(PROJECT_DIR, 'progress.log'), '--- PRD CHANGED: RESTARTING LOOP ---\n', 'utf8');
      lastHash = currentHash;
      if (watchMode) {
        i = 1;
      }
    }

    if (!watchMode && i > iterations) {
      console.log('⏸️  Max iterations reached.');
      break;
    }

    const phase = detectPhase();
    const iterId = `iter-${String(i).padStart(4, '0')}`;
    const iterDir = path.join(RUNS_DIR, iterId);
    ensureDir(iterDir);
    const tail = readTail(path.join(PROJECT_DIR, 'progress.log'), 200);
    fs.writeFileSync(path.join(iterDir, 'progress.tail.log'), tail, 'utf8');
    const latestLink = path.join(RUNS_DIR, 'latest');
    try {
      if (fileExists(latestLink)) fs.rmSync(latestLink, { recursive: true, force: true });
    } catch (_) {}
    try {
      fs.symlinkSync(iterDir, latestLink, 'junction');
    } catch (_) {}

    console.log('');
    console.log(`🔁 Loop ${i} (${phase} Phase)`);
    console.log(`   Logs: ${iterDir}`);

    const models = phase === 'PLAN' ? planModels : buildModels;
    let success = false;

    for (const model of models) {
      console.log(`   Using: ${model}`);
      const exitCode = await runAgent(model, phase, iterDir);
      const response = fs.readFileSync(path.join(iterDir, 'agent_response.txt'), 'utf8');

      if (/not supported|ModelNotFoundError|Make sure the model is enabled/i.test(response)) {
        console.log(`   ⚠️  Model ${model} not supported. Falling back...`);
        continue;
      }

      if (exitCode === 0 && response.trim().length > 0) {
        success = true;
        if (response.includes('<promise>COMPLETE</promise>')) {
          console.log('✅ Agent signaled completion.');
          if (!watchMode) {
            process.exit(0);
          }
          console.log('⏸️  Project Complete. Waiting for changes in prd.md...');
          while (md5File(path.join(PROJECT_DIR, 'prd.md')) === lastHash) {
            await new Promise((resolve) => setTimeout(resolve, 2000));
          }
          console.log('👀 Change detected! Resuming...');
          i = 1;
          break;
        }
        break;
      }

      console.log(`   ⚠️  Model ${model} failed (Exit: ${exitCode}). Falling back...`);
    }

    if (!success) {
      console.log('❌ All models failed this iteration.');
      await new Promise((resolve) => setTimeout(resolve, 2000));
    }

    lastHash = md5File(path.join(PROJECT_DIR, 'prd.md'));
    i += 1;
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
};

runLoop().catch((err) => {
  console.error('❌ Vibepup Windows runner failed.');
  console.error(String(err));
  process.exit(1);
});
