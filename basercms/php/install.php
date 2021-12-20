<?php
ini_set('display_errors', "Off");

$phpPath = "php";
$cake    = getcwd() . '/app/Console/cake.php';

$cmd  = "$phpPath -q $cake bc_manager install ";
// TODO
$cmd .= "http://site.url dbtype username password email --host hostname --database dbname --login dbuser --password dbpassword --prefix prefix_ --port portnumber --baseurl / --data nada-icons.default ";
$cmd .= "2>&1";

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
