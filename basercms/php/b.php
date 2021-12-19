<?php
$msg = "this is standard.\n";
$io = fopen('php://stdout', 'w');
fwrite($io, $msg);
fclose($io);

$msg = "this is error.\n";
$io = fopen('php://stderr', 'w');
fwrite($io, $msg);
fclose($io);

throw new Error("hoge");