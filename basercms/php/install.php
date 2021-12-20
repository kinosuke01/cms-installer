<?php
ini_set('display_errors', "Off");

function build_cmd($params = array()) {
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

$cmd = build_cmd($_POST);
$output = [];
$exit = 0;
if (!exec($cmd, $output, $exit)) {
  $output[] = 'unknown_error';
}
$res = [
  'exit_code' => $exit,
  'messages' => $output,
];
echo json_encode($res, JSON_UNESCAPED_UNICODE);
