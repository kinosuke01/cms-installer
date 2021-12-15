<?php
ini_set('display_errors', "On");
$argv = [
  (getcwd() . "/app/Console/cake.php"),
  "-working",
  (getcwd() . "/app"),
  "bc_manager",
  "install",
  "http://site.url",
  "dbtype",
  "username",
  "password",
  "email",
  "--host",
  "hostname",
  "--database",
  "dbname",
  "--login",
  "dbuser",
  "--password",
  "dbpassword",
  "--prefix",
  "prefix_",
  "--port",
  "portnumber",
  "--baseurl",
  "/",
  "--data",
  "nada-icons.default"
];
if (!defined('DS')) {
	define('DS', DIRECTORY_SEPARATOR);
}
$dispatcher = 'lib' . DS . 'Cake' . DS . 'Console' . DS . 'ShellDispatcher.php';

include($dispatcher);
unset($dispatcher);
return ShellDispatcher::run($argv);
