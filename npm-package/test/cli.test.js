const { test, describe, beforeEach } = require('node:test');
const assert = require('node:assert/strict');
const { execa } = require('execa');
const { createFixture } = require('fs-fixture');
const path = require('node:path');
const fs = require('node:fs');
const stripAnsi = require('strip-ansi');
const stripAnsiText = stripAnsi.default || stripAnsi;

const ralphBin = path.resolve('bin/ralph.js');

const writeOpencodeStub = (binDir) => {
  const stubPath = path.join(binDir, 'opencode');
  const content = `#!/usr/bin/env node
const fs = require('fs');
const path = require('path');

const args = process.argv.slice(2);
const command = args[0];

if (command === '--version') {
  console.log('opencode 0.0.0-test');
  process.exit(0);
}

if (command === 'models') {
  console.log('openai/gpt-5.2-codex');
  console.log('opencode/grok-code');
  process.exit(0);
}

if (command === 'auth') {
  console.log('auth ok');
  process.exit(0);
}

if (command === 'run') {
  console.log('<promise>COMPLETE</promise>');
  process.exit(0);
}

console.log('opencode stub');
process.exit(0);
`;

  fs.writeFileSync(stubPath, content, 'utf8');
  fs.chmodSync(stubPath, 0o755);

  return stubPath;
};

const runCli = async (args, { cwd, env }) => {
  const mergedEnv = {
    ...process.env,
    CI: 'true',
    GIT_TERMINAL_PROMPT: '0',
    GCM_INTERACTIVE: 'never',
    RALPH_MAX_TURN_SECONDS: '2',
    RALPH_NO_OUTPUT_SECONDS: '2',
    ...env
  };

  const result = await execa(process.execPath, [ralphBin, ...args], {
    cwd,
    env: mergedEnv,
    reject: false
  });

  return {
    ...result,
    stdout: stripAnsiText(result.stdout),
    stderr: stripAnsiText(result.stderr)
  };
};

describe('vibepup CLI end-to-end', () => {
  let fixture;

  beforeEach(async () => {
    fixture = await createFixture({
      'README.md': '# Fixture'
    });
  });

  test('runs default CLI flow with stubbed opencode', async () => {
    const binDir = path.join(fixture.path, 'bin');
    fs.mkdirSync(binDir, { recursive: true });
    writeOpencodeStub(binDir);

    const { stdout, exitCode } = await runCli(['1'], {
      cwd: fixture.path,
      env: {
        PATH: `${binDir}${path.delimiter}${process.env.PATH}`
      }
    });

    assert.equal(exitCode, 0);
    assert.ok(stdout.includes('🐾 Vibepup is waking up...'));
    assert.ok(stdout.includes('🔁 Loop 1'));

    const runsDir = path.join(fixture.path, '.ralph', 'runs');
    const iterDir = path.join(runsDir, 'iter-0001');
    assert.ok(fs.existsSync(iterDir));
    assert.ok(fs.existsSync(path.join(iterDir, 'agent_response.txt')));
  });

  test('runs free setup flow with stubbed opencode', async () => {
    const binDir = path.join(fixture.path, 'bin');
    fs.mkdirSync(binDir, { recursive: true });
    writeOpencodeStub(binDir);

    const { stdout, exitCode } = await runCli(['free'], {
      cwd: fixture.path,
      env: {
        PATH: `${binDir}${path.delimiter}${process.env.PATH}`
      }
    });

    assert.equal(exitCode, 0);
    assert.ok(stdout.includes('✨ Vibepup Free Setup'));
  });

  test('windows runner flow uses stubbed opencode', async () => {
    const binDir = path.join(fixture.path, 'bin');
    fs.mkdirSync(binDir, { recursive: true });
    writeOpencodeStub(binDir);

    const { stdout, exitCode } = await runCli(['--platform=windows', 'free'], {
      cwd: fixture.path,
      env: {
        PATH: `${binDir}${path.delimiter}${process.env.PATH}`
      }
    });

    assert.equal(exitCode, 0);
    assert.ok(stdout.includes('🐾 Vibepup is waking up...'));
    assert.ok(stdout.includes('✨ Vibepup Free Setup'));
  });
});
