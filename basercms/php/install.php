<?php
ini_set('display_errors', "Off");

$phpPath = "php";
$cake    = getcwd() . '/app/Console/cake.php';
$working = getcwd() . '/app';

$cmd  = "$phpPath -q $cake -woking $working bc_manager install ";
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
echo json_encode($res);
