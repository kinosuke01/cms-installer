package basercms

// TODO sync with install.php
const bcInstallScriptTemplate string = `<?php
const TOKEN      = 'TOKEN_PLACEHOLDER';
const EXPIRED_AT = 'EXPIRED_AT_PLACEHOLDER';
const PHP_PATH   = 'PHP_PATH_PLACEHOLDER';

ini_set('display_errors', "Off");

function validate_token($token, $opts = [])
{
	$correctToken = isset($opts['token']) ? $opts['token'] : TOKEN;
	$expiredAt = isset($opts['expired_at']) ? $opts['expired_at'] : EXPIRED_AT;
  $now = isset($opts['now']) ? $opts['now'] : strtotime('now');

  if ($correctToken !== $token) {
    return false;
  }
  if ((int)$expiredAt <= (int)$now) {
    return false;
  }
	return true;
}

function build_cmd($params = [], $opts = [])
{
  $phpPath = isset($opts['php_path']) ? $opts['php_path'] : PHP_PATH;
  $cake    = __DIR__ . '/app/Console/cake.php';

  $cmds = ["$phpPath -q $cake bc_manager install"];

  $reqKeys = ['siteurl', 'dbtype', 'siteuser', 'sitepassword', 'email'];
  foreach($reqKeys as $key) {
    $val = isset($params[$key]) ? $params[$key] : '';
    $cmds[] = escapeshellarg($val);
  }

  $optKeys = ['host', 'database', 'login', 'password', 'prefix', 'port', 'baseurl', 'data'];
  foreach($optKeys as $key) {
    if (isset($params[$key]) && $params[$key]) {
      $cmds[] = '--' . $key . ' ' . escapeshellarg($params[$key]);
    }
  }

  $cmds[] = '2>&1';

  $cmd = implode(' ', $cmds);
  return $cmd;
}

function res($exitCode = 0, $messages = array())
{
  $res = [
    'exit_code' => $exitCode,
    'messages' => $messages,
  ];
  return json_encode($res, JSON_UNESCAPED_UNICODE);
}

function main()
{
  if (!validate_token($_POST['token'])) {
    echo res(1, ['AUTH_ERROR']);
    return;
  }

  $cmd = build_cmd($_POST);
  $output = [];
  $exit = 0;
  if (!exec($cmd, $output, $exit)) {
    $output[] = 'EXEC_ERROR';
  }
  echo res($exit, $output);
}

if (isset($_POST['token'])) {
	main();
}
`
