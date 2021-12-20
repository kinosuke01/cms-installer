<?php
$pwd = __DIR__;
require_once("$pwd/install.php");

function errf($line, $desc, $case, $key, $want, $got) {
    return sprintf(
        "%s:%s %s:%s %s wrong. want=%s got=%s",
        __FILE__, $line, $desc, $case, $key, $want, $got
    );
}

$funcs = array(
    'validate_token' => function($desc) {
        $errs = array();
        $tt = array(
            array(
                'case' => 'invalid_token',
                'itoken' => 'abcde',
                'expired_at' => '1634292000', // '2021-10-15T19:00:00+09:00
                'token' => '12345',
                'now' => '1634288400',
                'expected_result' => false,
            ),
            array(
                'case' => 'expired',
                'itoken' => 'abcde',
                'expired_at' => '1634292000',
                'token' => 'abcde',
                'now' => '1634292000',
                'expected_result' => false,
            ),
            array(
                'case' => 'success',
                'itoken' => 'abcde',
                'expired_at' => '1634292000',
                'token' => 'abcde',
                'now' => '1634291999',
                'expected_result' => true,
            ),
        );
        foreach ($tt as $tc) {
            $result = validate_token($tc['token'], [
                'now' => $tc['now'],
                'token' => $tc['itoken'],
                'expired_at' => $tc['expired_at'],
            ]);

            if ($tc['expected_result'] !== $result) {
                $errs[] = errf(__LINE__, $desc, $tc['case'], 'error', $tc['expected_result'], $result);
            }
        }
        return $errs;
    },
);

$errs = array();
foreach ($funcs as $desc => $func) {
    $errs = array_merge($errs, $func($desc));
}
if (empty($errs)) {
    echo "TEST OK\n";
    exit(0);
} else {
    echo implode("\n", $errs);
    echo "\n";
    exit(1);
}
