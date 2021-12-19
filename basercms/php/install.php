<?php
$phpPath = "php";
$cmd     = getcwd() . '/app/Console/cake.php';
$working = getcwd() . '/app';

// sh ./app/Console/cake bc_manager install -app /Users/kinosuke01/Desktop/basercms/app "http://site.url" "dbtype" "username" "password" "email" --host "hostname" --database "dbname" --login "dbuser" --password "dbpassword" --prefix "prefix_" --port "portnumber" --baseurl "/" --data "nada-icons.default"

// exec php -q "$CONSOLE"/cake.php -working "$APP" "$@"

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
