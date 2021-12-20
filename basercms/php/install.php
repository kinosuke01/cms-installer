<?php
const TOKEN      = 'TOKEN_PLACEHOLDER';
const EXPIRED_AT = 'EXPIRED_AT_PLACEHOLDER';
const PHP_PATH   = 'PHP_PATH_PLACEHOLDER';

ini_set('display_errors', "Off");

function validate_token($token, $now = null)
{
	$correctToken = TOKEN;
	$expiredAt = EXPIRED_AT;

	if (!$now) {
	  $now = strtotime('now');
	}
  if ($correctToken !== $token) {
    return false;
  }
  if ((int)$expiredAt <= (int)$now) {
    return false;
  }
	return true;
}

function build_cmd($params = array())
{
  $phpPath = "php";
  $cake    = getcwd() . '/app/Console/cake.php';

  $cmds = ["$phpPath -q $cake bc_manager install"];

  $reqKeys = ['siteurl', 'dbtype', 'username', 'password'];
  foreach($reqKeys as $key) {
    $val = isset($params[$key]) ? $params[$key] : '';
    $cmds[] = escapeshellarg($val);
  }

  $optKeys = ['host', 'database', 'login', 'password', 'dbpassword', 'prefix', 'port', 'portnumber', 'baseurl', 'data'];
  foreach($optKeys as $key) {
    if (isset($params[$key])) {
      $cmds[] = '--' . $key . ' ' . escapeshellarg($params[$key]);
    }
  }

  $cmds[] = '2>&1';

  $cmd = implode(' ', $cmds);
  return $cmd;
}

function main()
{
  if (!validate_token($_POST['token'])) {
    // TODO
    return;
  }

  $cmd = build_cmd($_POST);
  $output = [];
  $exit = 0;
  if (!exec($cmd, $output, $exit)) {
    $output[] = 'UNKNOWN_ERROR';
  }
  $res = [
    'exit_code' => $exit,
    'messages' => $output,
  ];
  echo json_encode($res, JSON_UNESCAPED_UNICODE);
}

if (isset($_POST['token'])) {
	main();
}
